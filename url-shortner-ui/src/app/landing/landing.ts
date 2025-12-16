import { AfterViewInit, Component } from '@angular/core';
import { RouterLink } from '@angular/router';

import { Navbar } from '../navbar/navbar';

@Component({
  selector: 'app-landing',
  imports: [RouterLink, Navbar],
  templateUrl: './landing.html',
  styleUrl: './landing.scss',
})
export class Landing implements AfterViewInit {

  ngAfterViewInit() {
    const card = document.getElementById('card');

    if (!card) return;

    document.addEventListener('mousemove', (e) => {
      const x = (window.innerWidth / 2 - e.clientX) / 30;
      const y = (window.innerHeight / 2 - e.clientY) / 30;

      card.style.transform = `rotateX(${y}deg) rotateY(${-x}deg) translateZ(20px)`;
    });

    document.addEventListener('mouseleave', () => {
      card.style.transform = `rotateX(0deg) rotateY(0deg) translateZ(0)`;
    });
  }

}
