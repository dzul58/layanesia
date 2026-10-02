package controllers

import (
	"fmt"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService   services.UserService
	uploadService services.UploadService
}

func NewUserController(userService services.UserService, uploadService services.UploadService) *UserController {
	return &UserController{userService: userService, uploadService: uploadService}
}

func (ctrl *UserController) GetMe(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	user, err := ctrl.userService.GetProfile(userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"user": user})
}

func (ctrl *UserController) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.UpdateProfileDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	user, err := ctrl.userService.UpdateProfile(userID, req)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui!", "user": user})
}

func (ctrl *UserController) SwitchMode(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.SwitchModeDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	user, newToken, err := ctrl.userService.SwitchMode(userID, req.Mode)
	if err != nil {
		return err
	}
	modeLabel := "Pencari Kerja (SEEKER)"
	if user.ActiveMode == models.ModeEmployer {
		modeLabel = "Pemberi Kerja (EMPLOYER)"
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Mode akun berhasil diubah ke %s!", modeLabel),
		"token":   newToken,
		"user":    user,
	})
}

// VerifyKTP menerima pengajuan KTP. Status akhir tergantung konfigurasi
// (PENDING_REVIEW menunggu admin, atau VERIFIED jika KTP_AUTO_VERIFY di dev).
func (ctrl *UserController) VerifyKTP(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.VerifyKTPDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	user, err := ctrl.userService.SubmitKTP(userID, req)
	if err != nil {
		return err
	}
	msg := "Pengajuan verifikasi KTP diterima dan sedang ditinjau."
	if user.KTPStatus == models.KTPVerified {
		msg = "KTP Anda berhasil diverifikasi!"
	}
	return c.JSON(fiber.Map{"message": msg, "ktp_status": user.KTPStatus, "user": user})
}

func (ctrl *UserController) UploadResume(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.UploadResumeDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	user, err := ctrl.userService.UploadResume(userID, req)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "File Resume PDF berhasil disimpan!", "user": user})
}

// UploadFile menerima multipart `file` (+ opsional `purpose`=ktp|resume).
func (ctrl *UserController) UploadFile(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return apperrors.BadRequest("File tidak ditemukan dalam form upload (field 'file').")
	}
	purposeRaw := c.FormValue("purpose", c.Query("purpose"))
	purpose, err := services.ParseUploadPurpose(purposeRaw)
	if err != nil {
		return err
	}
	res, err := ctrl.uploadService.Upload(c.Context(), userID, purpose, fh)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message":   "File berhasil diunggah",
		"file_url":  res.URL,
		"mime_type": res.MIMEType,
		"bytes":     res.Bytes,
		"storage":   res.Storage,
	})
}

// --- Admin ---

func (ctrl *UserController) AdminListPendingKTP(c *fiber.Ctx) error {
	limit, offset := pagination(c, 50, 200)
	users, err := ctrl.userService.ListPendingKTP(limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": users, "count": len(users)})
}

func (ctrl *UserController) AdminReviewKTP(c *fiber.Ctx) error {
	targetID, err := paramUUID(c, "id", "ID pengguna")
	if err != nil {
		return err
	}
	var req services.ReviewKTPDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	user, err := ctrl.userService.ReviewKTP(targetID, req)
	if err != nil {
		return err
	}
	msg := "Pengajuan KTP ditolak."
	if req.Approve {
		msg = "KTP pengguna disetujui."
	}
	return c.JSON(fiber.Map{"message": msg, "user": user.Public(), "ktp_status": user.KTPStatus})
}
