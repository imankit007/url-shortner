import { HttpClient } from '@angular/common/http';
import { Injectable, computed, inject, signal } from '@angular/core';
import { Observable, tap } from 'rxjs';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthTokenResponse {
  access_token: string;
  token_type: string;
  expires_at: number;
  user_id: string;
  tenant_id: string;
  email: string;
}

const storageKey = 'url-shortener-auth-session';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly sessionSignal = signal<AuthTokenResponse | null>(this.loadSession());

  readonly session = this.sessionSignal.asReadonly();
  readonly user = computed(() => this.sessionSignal());
  readonly isAuthenticated = computed(() => this.sessionSignal() !== null);
  readonly accessToken = computed(() => this.sessionSignal()?.access_token ?? '');

  login(payload: LoginRequest): Observable<AuthTokenResponse> {
    return this.http
      .post<AuthTokenResponse>('http://localhost:8081/api/v1/auth/token', payload)
      .pipe(tap((session) => this.persistSession(session)));
  }

  logout(): void {
    this.sessionSignal.set(null);
    localStorage.removeItem(storageKey);
  }

  private persistSession(session: AuthTokenResponse): void {
    this.sessionSignal.set(session);
    localStorage.setItem(storageKey, JSON.stringify(session));
  }

  private loadSession(): AuthTokenResponse | null {
    const storedValue = localStorage.getItem(storageKey);
    if (!storedValue) {
      return null;
    }

    try {
      return JSON.parse(storedValue) as AuthTokenResponse;
    } catch {
      localStorage.removeItem(storageKey);
      return null;
    }
  }
}
