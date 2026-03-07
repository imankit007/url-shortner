import { HttpClient } from "@angular/common/http"
import { Injectable } from "@angular/core"
import { Observable } from "rxjs"


export interface ShortUrlRequest {
  links: Array<LinkRequest>
}

export interface LinkRequest {
  url: string
}


@Injectable({ providedIn: 'root' })
export class ShortUrlService {

  constructor(private http: HttpClient) { }

  createShortUrl(data: ShortUrlRequest): Observable<any> {
    return this.http.post<any>(
      'http://localhost:8080/shorten', data
    );
  }


}
