import { notFound } from 'next/navigation';
import LandingContent from './landing-content';

async function fetchLanding(id) {
  const base = process.env.NEXT_PUBLIC_API_BASE || 'https://adflow.up.railway.app/api/v1';
  const res = await fetch(`${base}/public/landing/${id}`, { cache: 'no-store' });
  if (!res.ok) {
    return null;
  }
  const body = await res.json();
  if (!body || !body.success || !body.data) {
    return null;
  }
  return body.data;
}

export default async function LandingPage({ params }) {
  const data = await fetchLanding(params.id);
  if (!data) {
    notFound();
  }

  const user = data.user || {};
  const landing = data.landing || {};
  return <LandingContent user={user} landing={landing} />;
}
