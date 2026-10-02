package models

// File ini berisi proyeksi (view) yang dipakai sebagai respons API.
// Relasi *User pada model sengaja diberi tag json:"-" sehingga secara default
// tidak ada data pribadi yang ikut terserialisasi; view di bawah menambahkan
// kembali informasi relasi dalam bentuk PublicUser yang aman.

// SkillPostingView = SkillPosting + profil publik pemilik.
type SkillPostingView struct {
	SkillPosting
	User *PublicUser `json:"user,omitempty"`
}

func (sp *SkillPosting) View() *SkillPostingView {
	if sp == nil {
		return nil
	}
	return &SkillPostingView{SkillPosting: *sp, User: sp.User.Public()}
}

func SkillPostingViews(items []SkillPosting) []SkillPostingView {
	out := make([]SkillPostingView, 0, len(items))
	for i := range items {
		out = append(out, *items[i].View())
	}
	return out
}

// JobPostingView = JobPosting + profil publik pemberi kerja.
type JobPostingView struct {
	JobPosting
	Employer *PublicUser `json:"employer,omitempty"`
}

func (jp *JobPosting) View() *JobPostingView {
	if jp == nil {
		return nil
	}
	return &JobPostingView{JobPosting: *jp, Employer: jp.Employer.Public()}
}

func JobPostingViews(items []JobPosting) []JobPostingView {
	out := make([]JobPostingView, 0, len(items))
	for i := range items {
		out = append(out, *items[i].View())
	}
	return out
}

// JobApplicationView dipakai untuk pelamar (melihat lamarannya sendiri) maupun
// pemberi kerja (meninjau pelamar). Applicant memuat resume tetapi bukan kontak.
type JobApplicationView struct {
	JobApplication
	JobPosting *JobPostingView `json:"job_posting,omitempty"`
	Applicant  *ApplicantView  `json:"applicant,omitempty"`
}

func (ja *JobApplication) View() *JobApplicationView {
	if ja == nil {
		return nil
	}
	return &JobApplicationView{
		JobApplication: *ja,
		JobPosting:     ja.JobPosting.View(),
		Applicant:      ja.Applicant.Applicant(),
	}
}

func JobApplicationViews(items []JobApplication) []JobApplicationView {
	out := make([]JobApplicationView, 0, len(items))
	for i := range items {
		out = append(out, *items[i].View())
	}
	return out
}

// JobOfferView = JobOffer + skill posting + profil publik kedua pihak.
type JobOfferView struct {
	JobOffer
	SkillPosting *SkillPostingView `json:"skill_posting,omitempty"`
	Employer     *PublicUser       `json:"employer,omitempty"`
	Worker       *PublicUser       `json:"worker,omitempty"`
}

func (jo *JobOffer) View() *JobOfferView {
	if jo == nil {
		return nil
	}
	return &JobOfferView{
		JobOffer:     *jo,
		SkillPosting: jo.SkillPosting.View(),
		Employer:     jo.Employer.Public(),
		Worker:       jo.Worker.Public(),
	}
}

func JobOfferViews(items []JobOffer) []JobOfferView {
	out := make([]JobOfferView, 0, len(items))
	for i := range items {
		out = append(out, *items[i].View())
	}
	return out
}

// RatingView = Rating + profil publik pemberi ulasan.
type RatingView struct {
	Rating
	Reviewer *PublicUser `json:"reviewer,omitempty"`
}

func (r *Rating) View() *RatingView {
	if r == nil {
		return nil
	}
	return &RatingView{Rating: *r, Reviewer: r.Reviewer.Public()}
}

func RatingViews(items []Rating) []RatingView {
	out := make([]RatingView, 0, len(items))
	for i := range items {
		out = append(out, *items[i].View())
	}
	return out
}
