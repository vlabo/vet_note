
import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { IonicModule } from '@ionic/angular';
import { PatientsService } from '../patients.service';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { ViewPatient } from '../types';

@Component({
  selector: 'app-edit-patient',
  templateUrl: './edit-patient.component.html',
  styleUrls: ['./edit-patient.component.scss'],
  standalone: true,
  imports: [CommonModule, FormsModule, IonicModule]
})
export class EditPatientComponent implements OnInit {
  patient: ViewPatient;
  types: String[] = [];
  newMode = false;

  constructor(
    private patientsService: PatientsService,
    private location: Location,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.patient = {
      id: "",
      type: "",
      name: "",
      gender: 'unknown',
      birthDate: "-", 
      chipId: "",
      weight: 0 /* float64 */,
      castrated: false,
      lastModified: "",
      note: "",
      owner: "",
      ownerPhone: "",
      procedures: [],
    };
  }

  ngOnInit(): void {
    this.route.url.subscribe(urlSegments => {
      var lastPathSegment = urlSegments[urlSegments.length - 1].path
      if (lastPathSegment === "new") {
        this.newMode = true;
      }
    });

    this.route.paramMap.subscribe(async params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      this.patientsService.getPatient(id).subscribe({
        next: patient => {
          console.log(patient);
          this.patient = patient;
        }
      })
    });
    this.patientsService.getPatientTypes().subscribe({
      next: types => {
        this.types = types;
      }
    });
  }

  async save() {
    this.patientsService.updatePatient(this.patient!).subscribe({
      next: patient => {
        if (this.newMode) {
          this.newMode = false;
          this.router.navigate(["/patient", patient.id], { replaceUrl: true })
          this.patient = patient;
        } else {
          this.location.back();
        }
      },
      error: error => {
        console.error('Error updating patient:', error);
      }
    }
    );
  }

  cancel(): void {
    this.location.back();
  }

  isValidDate(dateString: string): boolean {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  }

}