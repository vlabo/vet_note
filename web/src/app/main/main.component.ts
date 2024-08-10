import { ChangeDetectorRef, Component, OnDestroy, OnInit, SecurityContext, ViewChild, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { add, settingsSharp } from 'ionicons/icons';
import { faMicrochip, faUser, faMinus, faPhone } from '@fortawesome/free-solid-svg-icons';
import { IonSearchbar } from '@ionic/angular';
import Fuse from 'fuse.js';
import { ViewListPatient } from 'src/app/types';
import { PatientsService } from '../patients.service';
import { addIcons } from 'ionicons';
import { debounceTime, filter, Subject, Subscription } from 'rxjs';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss'],
  encapsulation: ViewEncapsulation.None // Disable view encapsulation
  })
export class MainComponent implements OnInit, OnDestroy {
  // Icons
  microchip = faMicrochip;
  user = faUser;
  minus = faMinus;
  phone = faPhone;

  private routerSubscription: Subscription | null = null;

  // searchbar reference used for filtering.
  @ViewChild('searchbar', { static: false }) searchbar!: IonSearchbar;

  patients: ViewListPatient[] = [];
  filteredPatients: any[] = [];

  fuse: Fuse<ViewListPatient>;

  private searchTerms = new Subject<string>();
  sanitizer: DomSanitizer
  constructor(private router: Router, private route: ActivatedRoute, private patientsService: PatientsService, private s: DomSanitizer, private cdr: ChangeDetectorRef) {
    this.fuse = new Fuse(this.patients, {});
    this.sanitizer = s;
  }

  async ngOnInit() {
    // Initialization.
    addIcons({ "add": add, "settings": settingsSharp });

    this.subscribeToPatientList();
    this.patientsService.triggerPatientListReload();

    this.searchTerms.pipe(
      debounceTime(300) // Wait for 300ms pause in events
    ).subscribe(query => {
      this.filterList(query)
    });
  }

  ngOnDestroy(): void {
    if (this.routerSubscription) {
      this.routerSubscription.unsubscribe();
    }
  }

  subscribeToPatientList() {
    this.patientsService.getPatientListObservable().subscribe({
      next: patients => {
        this.patients = patients;
        this.filteredPatients = this.patients;
        this.fuse = new Fuse(this.patients, {
          keys: ['name', 'owner', 'chipId', "phone"],
          includeMatches: true,
        });

        // the query parameter from the URL and filter the list based on it.
        let searchQuery = this.route.snapshot.queryParamMap.get('query');
        this.filterList(searchQuery);
        this.searchbar.value = searchQuery;
      }
    });
  }

  goToPatient(patientId: string): void {
    this.router.navigate(['/patient', patientId]);
  }

  newPatient(): void {
    this.router.navigate(['/patient/new']);
  }

  openSettings(): void {
    this.router.navigate(['/settings']);
  }

  // onSearch is called when the search input changes. It filters the list of patients based on the search query.
  onSearch(event: any): void {
    this.searchTerms.next(event.target.value);
  }

  // filterList calls the filter function of the fuse.js library to filter the patient list.
  filterList(query: string | null): void {
    // this.router.navigate([], { queryParams: { query: query }, queryParamsHandling: 'merge', replaceUrl: true });
    if(query) {
      const result = this.fuse.search(query);
      this.filteredPatients = result.map((res): ViewListPatient  => {
        let copy : any = {
          id: res.item.id,
          type: this.escapeHtml(res.item.type),
          name: this.escapeHtml(res.item.name),
          owner: this.escapeHtml(res.item.owner),
        }
        res.matches?.forEach( match => {
          // @ts-ignore
          let text = this.escapeHtml(res.item[match.key as keyof ViewListPatient]); 
          if(!text) {
            return;
          }
          let indices = match.indices.slice().reverse();
          indices.forEach(([start, end]) => {
            text = text!.substring(0, start) + '<span class="bold-blue">' + text!.substring(start, end + 1) + '</span>' + text!.substring(end + 1);
          });
          // @ts-ignore
          copy[match.key] = this.sanitizer.bypassSecurityTrustHtml(text);
        });

        return copy;
      });
    } else {
      this.filteredPatients = this.patients.map((p): ViewListPatient  => {
        let copy : any = {
          id: p.id,
          type: this.escapeHtml(p.type),
          name: this.escapeHtml(p.name),
          owner: this.escapeHtml(p.owner),
        }
        return copy;
      });
    }
  }

  escapeHtml(input: string): string {
    if (!input) {
      return '';
    }
    return input.replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }
}
