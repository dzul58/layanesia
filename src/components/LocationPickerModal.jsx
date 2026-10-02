import React, { useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { MapPin, Search, X, Check } from 'lucide-react';
import { api } from '../api';

// Custom Marker Icon for Leaflet
const customPin = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41]
});

function MapEventsHandler({ onPositionSelect }) {
  useMapEvents({
    click(e) {
      onPositionSelect(e.latlng.lat, e.latlng.lng);
    },
  });
  return null;
}

export default function LocationPickerModal({ isOpen, onClose, onConfirmLocation, initialLat = -6.2435, initialLng = 106.8021 }) {
  const [position, setPosition] = useState({ lat: initialLat, lng: initialLng });
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  const [formattedAddress, setFormattedAddress] = useState('Kebayoran Baru, Jakarta Selatan, Indonesia');
  const [osmPlaceId, setOsmPlaceId] = useState('osm_31740123');

  if (!isOpen) return null;

  const handleSearchNominatim = async (e) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;
    setIsSearching(true);
    try {
      const data = await api(`/api/v1/locations/nominatim/search?q=${encodeURIComponent(searchQuery)}`);
      setSearchResults(data.data || []);
    } catch (err) {
      console.error('Nominatim search error:', err);
      setSearchResults([]);
    } finally {
      setIsSearching(false);
    }
  };

  const handleSelectSearchResult = (result) => {
    const lat = parseFloat(result.lat);
    const lng = parseFloat(result.lon);
    setPosition({ lat, lng });
    setFormattedAddress(result.display_name);
    setOsmPlaceId(String(result.place_id || 'osm_' + Date.now()));
    setSearchResults([]);
  };

  const handlePositionSelect = (lat, lng) => {
    setPosition({ lat, lng });
    setFormattedAddress(`Titik Koordinat: ${lat.toFixed(6)}, ${lng.toFixed(6)} (OpenStreetMap)`);
    setOsmPlaceId(`osm_${lat.toFixed(4)}_${lng.toFixed(4)}`);
  };

  const handleSave = () => {
    onConfirmLocation({
      latitude: position.lat,
      longitude: position.lng,
      formattedAddress: formattedAddress,
      osmPlaceId: osmPlaceId,
    });
    onClose();
  };

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      zIndex: 1000,
      background: 'rgba(0, 0, 0, 0.7)',
      backdropFilter: 'blur(8px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <div className="glass-panel" style={{ width: '100%', maxWidth: '720px', padding: '24px', background: 'var(--bg-secondary)', position: 'relative', border: '1px solid var(--border-color)' }}>
        {/* Close Button */}
        <button onClick={onClose} style={{ position: 'absolute', top: '20px', right: '20px', background: 'transparent', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}>
          <X size={24} />
        </button>

        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '16px' }}>
          <MapPin color="#10b981" size={24} />
          <h2 style={{ fontSize: '1.25rem', color: 'var(--text-main)' }}>Pilih Titik Lokasi (OpenStreetMap & Leaflet.js)</h2>
        </div>

        {/* Nominatim Search Input */}
        <form onSubmit={handleSearchNominatim} style={{ display: 'flex', gap: '10px', marginBottom: '16px' }}>
          <div style={{ flex: 1, position: 'relative' }}>
            <input 
              type="text" 
              placeholder="Cari lokasi via OpenStreetMap Nominatim (misal: Kebayoran Baru, Jakarta)..."
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              style={{ width: '100%', padding: '10px 14px', borderRadius: '8px', background: 'var(--input-bg)', border: '1px solid var(--border-color)', color: 'var(--text-main)', fontSize: '0.9rem' }}
            />
          </div>
          <button type="submit" className="btn btn-primary" style={{ height: '40px' }} disabled={isSearching}>
            <Search size={16} /> {isSearching ? 'Cari...' : 'Cari Alamat'}
          </button>
        </form>

        {/* Search Results Dropdown */}
        {searchResults.length > 0 && (
          <div style={{ background: 'var(--bg-glass)', borderRadius: '8px', padding: '8px', marginBottom: '16px', maxHeight: '160px', overflowY: 'auto', border: '1px solid var(--border-color)' }}>
            {searchResults.map((item, idx) => (
              <div 
                key={idx} 
                onClick={() => handleSelectSearchResult(item)}
                style={{ padding: '8px 12px', borderRadius: '6px', cursor: 'pointer', fontSize: '0.85rem', borderBottom: '1px solid var(--border-color)', color: 'var(--text-main)' }}
              >
                📍 {item.display_name}
              </div>
            ))}
          </div>
        )}

        {/* Leaflet Map Box */}
        <div style={{ height: '320px', borderRadius: '12px', overflow: 'hidden', border: '1px solid var(--border-color)', marginBottom: '16px' }}>
          <MapContainer center={[position.lat, position.lng]} zoom={13} style={{ height: '100%', width: '100%' }}>
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <Marker position={[position.lat, position.lng]} icon={customPin}>
              <Popup>
                📍 <strong>Titik Terpilih</strong><br />
                {formattedAddress}
              </Popup>
            </Marker>
            <MapEventsHandler onPositionSelect={handlePositionSelect} />
          </MapContainer>
        </div>

        {/* Selected Coordinates & Confirmation */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px', background: 'var(--bg-glass)', padding: '12px 16px', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
          <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
            <div style={{ color: '#10b981', fontWeight: 600 }}>Lat: {position.lat.toFixed(6)}, Lng: {position.lng.toFixed(6)}</div>
            <div style={{ fontSize: '0.8rem', textOverflow: 'ellipsis', overflow: 'hidden', maxWidth: '400px', whiteSpace: 'nowrap', color: 'var(--text-main)' }}>{formattedAddress}</div>
          </div>

          <div style={{ display: 'flex', gap: '10px' }}>
            <button onClick={onClose} className="btn btn-outline" style={{ fontSize: '0.85rem' }}>Batal</button>
            <button onClick={handleSave} className="btn btn-primary" style={{ fontSize: '0.85rem' }}>
              <Check size={16} /> Simpan Titik Lokasi
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
