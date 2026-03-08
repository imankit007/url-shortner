import { Component } from '@angular/core';
 import { NgxShineBorderComponent } from '@omnedia/ngx-shine-border';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-navbar',
  imports: [ButtonModule, NgxShineBorderComponent],
  templateUrl: './navbar.html',
  styleUrl: './navbar.scss',
})
export class Navbar {

}
