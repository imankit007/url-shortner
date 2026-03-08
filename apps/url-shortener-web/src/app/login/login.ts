import { CommonModule } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';
import { AuthService } from '../core/auth.service';
import { Navbar } from '../navbar/navbar';

@Component({
  selector: 'app-login',
  imports: [CommonModule, ReactiveFormsModule, RouterLink, Navbar],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  private readonly formBuilder = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal('');

  protected readonly loginForm = this.formBuilder.nonNullable.group({
    email: ['owner@tenant1.local', [Validators.required, Validators.email]],
    password: ['tenant1-password', Validators.required],
  });

  protected submit(): void {
    if (this.loginForm.invalid || this.isSubmitting()) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    this.errorMessage.set('');

    this.authService
      .login(this.loginForm.getRawValue())
      .pipe(finalize(() => this.isSubmitting.set(false)))
      .subscribe({
        next: () => {
          void this.router.navigate(['/dashboard']);
        },
        error: () => {
          this.errorMessage.set('Login failed. Use one of the demo tenant accounts.');
        },
      });
  }

  protected useTenantOne(): void {
    this.loginForm.setValue({
      email: 'owner@tenant1.local',
      password: 'tenant1-password',
    });
  }

  protected useTenantTwo(): void {
    this.loginForm.setValue({
      email: 'owner@tenant2.local',
      password: 'tenant2-password',
    });
  }
}
