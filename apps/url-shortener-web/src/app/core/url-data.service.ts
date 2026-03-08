import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

export interface ShortUrlRequest {
  links: Array<LinkRequest>;
}

export interface LinkRequest {
  url: string;
}

export interface ShortUrlResponse {
  long_url: string;
  short_url: string;
}

export interface UrlEntry {
  url: string;
  code: number;
  tenant_id: string;
  created_by_user_id: string;
}

@Injectable({ providedIn: 'root' })
export class UrlDataService {
  private readonly http = inject(HttpClient);
  private readonly apiBaseUrl = 'http://localhost:8080/api/v1';

  createShortUrls(data: ShortUrlRequest): Observable<Array<ShortUrlResponse>> {
    return this.http.post<Array<ShortUrlResponse>>(`${this.apiBaseUrl}/urls/shorten`, data);
  }

  listTenantUrls(): Observable<Array<UrlEntry>> {
    return this.http.get<Array<UrlEntry>>(`${this.apiBaseUrl}/urls`);
  }
}
