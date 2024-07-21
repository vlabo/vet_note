import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, OnInit, Renderer2, ViewChild, ViewEncapsulation } from '@angular/core';
import { IonSearchbar, IonicModule } from '@ionic/angular';
import { add } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { Patient, PatientsService } from '../patients.service';
import { ActivatedRoute, Router } from '@angular/router';
import Fuse, { FuseResult } from 'fuse.js';
import { faMicrochip, faUser, faMinus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  standalone: true,
  imports: [IonicModule, CommonModule, FontAwesomeModule],
  styleUrls: ['./main.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class MainComponent implements OnInit, AfterViewInit {
  // Icons
  microchip = faMicrochip;
  user = faUser;
  minus = faMinus;

  @ViewChild('searchbar', { static: false }) searchbar!: IonSearchbar;

  patients: Patient[] = [];

  filteredPatients: Patient[] = [];
  fuse: Fuse<Patient>;

  searchQuery: string | null = "";
  searchResult: FuseResult<Patient>[] | undefined;

  constructor(private router: Router, private route: ActivatedRoute, private patientsService: PatientsService) {
    this.fuse = new Fuse(this.patients, {});
  }

  ngOnInit() {
    addIcons({ "add": add });
    this.patients = this.patientsService.getPatients()
    this.filteredPatients = this.patients;
    this.fuse = new Fuse(this.patients, {
      keys: ['Name', 'Owner.Name', 'IdNumber'],
      includeMatches: true,
    });

    this.searchQuery = this.route.snapshot.queryParamMap.get('query');
    if (this.searchQuery) {
      this.filterList(this.searchQuery);
    }
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

  onSearch(event: any): void {
    const query = event.target.value;
    this.searchQuery = query;
    if (!query) {
      this.router.navigate([], {  relativeTo: this.route, replaceUrl: true });
      this.filteredPatients = this.patients;
      this.searchResult = [];
    } else {
      this.filterList(query)
    }
  }

  filterList(query: string): void {
    this.router.navigate([], { queryParams: { query: query }, queryParamsHandling: 'merge', replaceUrl: true});
    const result = this.fuse.search(query);
    this.searchResult = result;
    this.filteredPatients = result.map(res => res.item);
  }

  highlight(p: Patient, key: string, includeEmpty: boolean = true): string {
    let highlightedText: string = p[key as keyof Patient] as string;
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
