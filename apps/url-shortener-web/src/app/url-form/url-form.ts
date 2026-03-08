import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

import { CommonModule } from '@angular/common';
import { ShortUrlRequest, ShortUrlService } from './short-url.service';


@Component({
  selector: 'app-url-form',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './url-form.html',
  styleUrl: './url-form.scss',
})
export class UrlForm implements OnInit {


  shortUrlForm;

  constructor(private fb: FormBuilder, private shortUrlService: ShortUrlService) {
    this.shortUrlForm =
      this.fb.group({
        links: this.fb.array([this.createLink()])
      });
  }


  ngOnInit() {
    this.shortUrlForm.get('longUrl')?.valueChanges.subscribe(value => {
      console.log("Form changed : ", value);
    })
  }


  onSubmit() {
    if (this.shortUrlForm.valid) {
      console.log(JSON.stringify(this.shortUrlForm.value, null, 2));

      var req: ShortUrlRequest = {
        links: []
      };


      this.shortUrlForm.value.links?.forEach(
        link => {
          req.links.push({
            url: link.url
          })
        }
      )

      var response = this.shortUrlService.createShortUrl(req)

      response.subscribe({

        next: value => console.log(value),
        error: err => console.error(err),
        complete: () => console.log("Complete")
      });


    }
  }

  addLink() {
    this.links.push(this.createLink())
  }


  removeLink(index: number) {
    this.links.removeAt(index)
  }

  createLink(): FormGroup {
    return this.fb.group({
      url: ["", Validators.required]
    });
  }

  get links(): FormArray {
    return this.shortUrlForm.get('links') as FormArray;
  }

}
