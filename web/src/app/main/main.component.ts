import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, OnInit, Renderer2, ViewChild, ViewEncapsulation } from '@angular/core';
// import { IonSearchbar, IonicModule } from '@ionic/angular';
import { add, settingsSharp } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { PatientsService } from '../patients.service';
import { ActivatedRoute, Router } from '@angular/router';
import Fuse, { FuseResult } from 'fuse.js';
import { faMicrochip, faUser, faMinus, faPhone } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ViewListPatient } from '../types';
import { IonButton, IonCol, IonContent, IonFab, IonFabButton, IonGrid, IonHeader, IonIcon, IonItem, IonLabel, IonRow, IonSearchbar, IonTitle, IonToolbar } from '@ionic/angular/standalone';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  standalone: true,
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonLabel, IonSearchbar, IonCol, IonRow, IonGrid, IonIcon, IonItem, IonFabButton, IonFab, IonButton, CommonModule, FontAwesomeModule],
  styleUrls: ['./main.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class MainComponent implements OnInit, AfterViewInit {
  // Icons
  microchip = faMicrochip;
  user = faUser;
  minus = faMinus;
  phone = faPhone;

  // searchbar reference used for filtering.
  @ViewChild('searchbar', { static: false }) searchbar!: IonSearchbar;

  patients: ViewListPatient[] = [];
  filteredPatients: ViewListPatient[] = [];

  fuse: Fuse<ViewListPatient>;

  searchQuery: string | null = "";
  searchResult: FuseResult<ViewListPatient>[] | undefined;

  constructor(private router: Router, private route: ActivatedRoute, private patientsService: PatientsService) {
    this.fuse = new Fuse(this.patients, {});
  }

  async ngOnInit() {
    // Initialization.
    addIcons({ "add": add, "settings": settingsSharp });

    // TODO: if new patients come from the server this need to be updated. 
    this.patientsService.getPatientList().subscribe({
      next: patients => {
        this.patients = patients;
        this.filteredPatients = this.patients;
        this.fuse = new Fuse(this.patients, {
          keys: ['name', 'owner', 'chip_id', "phone"],
          includeMatches: true,
        });

        // the query parameter from the URL and filter the list based on it.
        this.searchQuery = this.route.snapshot.queryParamMap.get('query');
        if (this.searchQuery) {
          this.filterList(this.searchQuery);
        }
      }
    })
  }

  ngAfterViewInit(): void {
    if (this.searchbar && this.searchQuery) {
      this.searchbar.value = this.searchQuery;
    }
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
    const query = event.target.value;
    this.searchQuery = query;
    if (!query) {
      this.router.navigate([], { relativeTo: this.route, replaceUrl: true });
      this.filteredPatients = this.patients;
      this.searchResult = [];
    } else {
      this.filterList(query)
    }
  }

  // filterList calls the filter function of the fuse.js library to filter the patient list.
  filterList(query: string): void {
    this.router.navigate([], { queryParams: { query: query }, queryParamsHandling: 'merge', replaceUrl: true });
    const result = this.fuse.search(query);
    this.searchResult = result;
    this.filteredPatients = result.map(res => res.item);
  }

  // highlight checks the search result and adds blue bold tag to all the matched characters of they key.
  highlight(p: ViewListPatient, key: string, includeEmpty: boolean = true): string {
    let highlightedText: string = p[key as keyof ViewListPatient] as string;
    var hasMatch = false;
    if (this.searchResult) {
      this.searchResult.forEach(m => {
        if (m.item !== p) {
          return;
        }
        m.matches?.forEach(match => {
          if (match.key !== key) {
            return;
          }
          hasMatch = true;

          var indices = match.indices.slice().reverse();
          indices.forEach(([start, end]) => {
            highlightedText = highlightedText.substring(0, start) + '<span class="bold-blue">' + highlightedText.substring(start, end + 1) + '</span>' + highlightedText.substring(end + 1);
          });
        });
      });
    }
    if (!hasMatch && !includeEmpty) {
      return "";
    }
    return highlightedText;
  }
}
