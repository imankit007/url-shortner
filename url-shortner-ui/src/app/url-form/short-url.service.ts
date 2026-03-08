import { HttpClient } from "@angular/common/http"
import { Injectable } from "@angular/core"
import { Observable } from "rxjs"


export interface ShortUrlRequest {
  links: Array<LinkRequest>
}

export interface LinkRequest {
  url: string
}

export interface ShortUrlResponse {
  long_url: string,
  short_url: string
}

@Injectable({ providedIn: 'root' })
export class ShortUrlService {

  constructor(private http: HttpClient) { }

  createShortUrl(data: ShortUrlRequest): Observable<Array<ShortUrlResponse>> {
    return this.http.post<any>(
      'http://localhost:8080/api/v1/urls/shorten', data
    );
  }


}
