import React, { useState, useEffect } from 'react';
import { api } from '../api';

export default function LocationSelector({ onLocationChange, value, onChange }) {
  const notify = onChange || onLocationChange;
  const [provinces, setProvinces] = useState([]);
  const [cities, setCities] = useState([]);
  const [districts, setDistricts] = useState([]);

  const [selectedProvince, setSelectedProvince] = useState(value?.province || 'DKI Jakarta');
  const [selectedCity, setSelectedCity] = useState(value?.city || 'Jakarta Selatan');
  const [selectedDistrict, setSelectedDistrict] = useState(value?.district || 'Kebayoran Baru');

  useEffect(() => {
    if (!value) return;
    if (value.province && value.province !== selectedProvince) setSelectedProvince(value.province);
    if (value.city && value.city !== selectedCity) setSelectedCity(value.city);
    if (value.district && value.district !== selectedDistrict) setSelectedDistrict(value.district);
  }, [value?.province, value?.city, value?.district]);

  useEffect(() => {
    const fetchProvinces = async () => {
      try {
        const data = await api('/api/v1/locations/provinces');
        if (data.data) setProvinces(data.data);
      } catch {
        setProvinces(['DKI Jakarta', 'Jawa Barat', 'Jawa Timur', 'Banten', 'Jawa Tengah']);
      }
    };
    fetchProvinces();
  }, []);

  useEffect(() => {
    if (!selectedProvince) return;
    const fetchCities = async () => {
      try {
        const data = await api(`/api/v1/locations/cities?province=${encodeURIComponent(selectedProvince)}`);
        if (data.data) {
          setCities(data.data);
          if (data.data.length > 0 && !data.data.includes(selectedCity)) {
            setSelectedCity(data.data[0]);
          }
        }
      } catch {
        setCities(['Jakarta Selatan', 'Jakarta Pusat', 'Jakarta Barat', 'Jakarta Timur', 'Jakarta Utara']);
      }
    };
    fetchCities();
  }, [selectedProvince]);

  useEffect(() => {
    if (!selectedCity) return;
    const fetchDistricts = async () => {
      try {
        const data = await api(`/api/v1/locations/districts?city=${encodeURIComponent(selectedCity)}`);
        if (data.data) {
          const districtNames = data.data.map((d) => (typeof d === 'string' ? d : d.district));
          setDistricts(districtNames);
          if (districtNames.length > 0 && !districtNames.includes(selectedDistrict)) {
            setSelectedDistrict(districtNames[0]);
          }
        }
      } catch {
        setDistricts(['Kebayoran Baru', 'Cilandak', 'Kuningan', 'Pondok Indah', 'Tanah Abang']);
      }
    };
    fetchDistricts();
  }, [selectedCity]);

  useEffect(() => {
    if (notify) {
      notify({
        province: selectedProvince,
        city: selectedCity,
        district: selectedDistrict,
      });
    }
  }, [selectedProvince, selectedCity, selectedDistrict]);

  const selectStyle = {
    width: '100%',
    padding: '10px 14px',
    borderRadius: '8px',
    background: 'var(--input-bg)',
    border: '1px solid var(--border-color)',
    color: 'var(--text-main)',
    fontSize: '0.9rem',
    outline: 'none',
  };

  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '14px', width: '100%' }}>
      <div>
        <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '4px', textTransform: 'uppercase', fontWeight: 600 }}>
          Provinsi
        </label>
        <select
          value={selectedProvince}
          onChange={(e) => setSelectedProvince(e.target.value)}
          style={selectStyle}
        >
          {provinces.map((p) => (
            <option key={p} value={p} style={{ background: 'var(--bg-secondary)', color: 'var(--text-main)' }}>{p}</option>
          ))}
        </select>
      </div>

      <div>
        <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '4px', textTransform: 'uppercase', fontWeight: 600 }}>
          Kota / Kabupaten
        </label>
        <select
          value={selectedCity}
          onChange={(e) => setSelectedCity(e.target.value)}
          style={selectStyle}
        >
          {cities.map((c) => (
            <option key={c} value={c} style={{ background: 'var(--bg-secondary)', color: 'var(--text-main)' }}>{c}</option>
          ))}
        </select>
      </div>

      <div>
        <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '4px', textTransform: 'uppercase', fontWeight: 600 }}>
          Kecamatan
        </label>
        <select
          value={selectedDistrict}
          onChange={(e) => setSelectedDistrict(e.target.value)}
          style={selectStyle}
        >
          {districts.map((d) => (
            <option key={d} value={d} style={{ background: 'var(--bg-secondary)', color: 'var(--text-main)' }}>{d}</option>
          ))}
        </select>
      </div>
    </div>
  );
}
