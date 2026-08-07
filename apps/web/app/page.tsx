'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';

export default function Root() {
  const router = useRouter();
  useEffect(() => {
    const token = api.auth.getAccessToken();
    router.replace(token ? '/dashboard' : '/login');
  }, [router]);
  return null;
}
