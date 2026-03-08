import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormArray, FormGroup, NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { finalize } from 'rxjs';
import { AuthService } from '../core/auth.service';
import { ShortUrlResponse, UrlDataService, UrlEntry } from '../core/url-data.service';
import { Navbar } from '../navbar/navbar';

type DomainMetric = { domain: string; count: number; width: number };

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, ReactiveFormsModule, Navbar],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss',
})
export class Dashboard {
  private readonly formBuilder = inject(NonNullableFormBuilder);
  private readonly urlDataService = inject(UrlDataService);
  protected readonly authService = inject(AuthService);

  protected readonly isSubmitting = signal(false);
  protected readonly isRefreshing = signal(false);
  protected readonly errorMessage = signal('');
  protected readonly latestResponses = signal<Array<ShortUrlResponse>>([]);
  protected readonly urlEntries = signal<Array<UrlEntry>>([]);

  protected readonly createUrlForm = this.formBuilder.group({
    links: this.formBuilder.array([this.createLinkGroup()]),
  });

  protected readonly totalLinks = computed(() => this.urlEntries().length);
  protected readonly uniqueDomains = computed(
    () => new Set(this.urlEntries().map((entry) => this.extractDomain(entry.url))).size,
  );
  protected readonly radialValue = computed(() => {
    const totalCount = this.totalLinks();
    if (totalCount === 0) {
      return 0;
    }

    return Math.min(100, Math.round((this.uniqueDomains() / totalCount) * 100));
  });
  protected readonly topDomains = computed<Array<DomainMetric>>(() => {
    const domainCounts = new Map<string, number>();
    for (const entry of this.urlEntries()) {
      const domain = this.extractDomain(entry.url);
      domainCounts.set(domain, (domainCounts.get(domain) ?? 0) + 1);
    }

    const rankedDomains = Array.from(domainCounts.entries())
      .sort((left, right) => right[1] - left[1])
      .slice(0, 5);

    const maxCount = rankedDomains[0]?.[1] ?? 1;
    return rankedDomains.map(([domain, count]) => ({
      domain,
      count,
      width: Math.max(18, Math.round((count / maxCount) * 100)),
    }));
  });

  constructor() {
    this.refreshDashboard();
  }

  protected submit(): void {
    if (this.createUrlForm.invalid || this.isSubmitting()) {
      this.createUrlForm.markAllAsTouched();
      return;
    }

    const request = {
      links: this.createUrlForm.getRawValue().links.map((link) => ({
        url: link.url,
      })),
    };

    this.isSubmitting.set(true);
    this.errorMessage.set('');

    this.urlDataService
      .createShortUrls(request)
      .pipe(finalize(() => this.isSubmitting.set(false)))
      .subscribe({
        next: (responses) => {
          this.latestResponses.set(responses);
          this.resetForm();
          this.refreshDashboard();
        },
        error: () => {
          this.errorMessage.set('Failed to create short URLs. Check your session and the API services.');
        },
      });
  }

  protected refreshDashboard(): void {
    this.isRefreshing.set(true);
    this.errorMessage.set('');

    this.urlDataService
      .listTenantUrls()
      .pipe(finalize(() => this.isRefreshing.set(false)))
      .subscribe({
        next: (entries) => {
          this.urlEntries.set(entries);
        },
        error: () => {
          this.errorMessage.set('Failed to load tenant analytics.');
        },
      });
  }

  protected addLink(): void {
    this.links.push(this.createLinkGroup());
  }

  protected removeLink(index: number): void {
    if (this.links.length === 1) {
      this.links.at(0).setValue({ url: '' });
      return;
    }

    this.links.removeAt(index);
  }

  protected copyToClipboard(value: string): void {
    void navigator.clipboard.writeText(value);
  }

  protected get links(): FormArray<FormGroup> {
    return this.createUrlForm.controls.links as FormArray<FormGroup>;
  }

  private createLinkGroup() {
    return this.formBuilder.group({
      url: ['', Validators.required],
    });
  }

  private resetForm(): void {
    this.createUrlForm.setControl('links', this.formBuilder.array([this.createLinkGroup()]));
  }

  private extractDomain(url: string): string {
    const normalizedUrl = url.startsWith('http://') || url.startsWith('https://') ? url : `https://${url}`;

    try {
      return new URL(normalizedUrl).hostname.replace(/^www\./, '');
    } catch {
      return 'invalid-domain';
    }
  }
}
