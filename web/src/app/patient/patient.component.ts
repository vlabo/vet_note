import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PatientsService } from '../patients.service';
import { AlertController, IonModal } from '@ionic/angular';

import { faCakeCandles, faUser, faCalendarDays, faMarsAndVenus, faMicrochip, faMars, faVenus, faPhone, faWeightHanging, faCheck, faX } from '@fortawesome/free-solid-svg-icons';

import { add, create, checkmark, close, arrowBack, chevronForward, paw, person, calendar, heartHalf } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { Location } from '@angular/common';
import { ViewPatient, ViewProcedure } from '../types';

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss'],
})
export class PatientComponent implements OnInit, OnDestroy {

  // Icons
  cake = faCakeCandles;
  user = faUser;
  calendar = faCalendarDays;
  gender = faMarsAndVenus;
  microchip = faMicrochip;
  mars = faMars;
  venus = faVenus;
  phone = faPhone;
  weight = faWeightHanging;
  check = faCheck;
  xIcon = faX;

  patient: ViewPatient | undefined = undefined;
  procedures: ViewProcedure[] = [];
  originalPatient: any;
  isViewProcedure = false;

  selectedProcedure: ViewProcedure | null = null;

  @ViewChild('procedureDetailsModal', { static: true }) modal!: IonModal;

  constructor(private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private patientService: PatientsService,
    private alertController: AlertController,
  ) {
    addIcons({
      "add": add,
      "create": create,
      "checkmark": checkmark,
      "close": close,
      "arrow-back": arrowBack,
      "chevron-forward": chevronForward,
      "paw": paw,
      "person": person,
      "calendar": calendar,
      "heart-half": heartHalf,
    });
  }

  ngOnInit() {
    this.route.paramMap.subscribe(async params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      this.patientService.getPatient(id).subscribe({
        next: patient => {
          this.patient = patient;
          console.log(this.patient);
        }
      })
    });
  }

  ngOnDestroy() {
    this.isViewProcedure = false;
    this.updatePatient();
  }

  openEdit() {
    this.router.navigate(["/patient", this.patient?.id, "edit"]);
  }

  goBack(): void {
    this.location.back();
  }

  setOpen(isOpen: boolean) {
    this.isViewProcedure = isOpen;
  }

  openViewProcedure(procedure: ViewProcedure) {
    if (this.patient) {
      this.router.navigate(['/procedure', procedure.id]);
    }
  }

  createProcedure(): void {
    if (this.patient) {
      this.router.navigate(['/procedure', this.patient.id, 'new']);
      this.isViewProcedure = false;
    }
  }

  updatePatient() {
    if (this.patient) {
      this.patientService.updatePatient(this.patient).subscribe({});
    }
  }

  isValidDate(dateString: string): boolean {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  }

  getAge(): string {
    let date = new Date(this.patient!.birthDate);
    // Check if the birth date is valid
    if (isNaN(date.getTime())) {
      return "-";
    }
    let now = new Date();
    let years = now.getFullYear() - date.getFullYear();
    let months = now.getMonth() - date.getMonth();

    // Adjust if the month difference is negative
    if (months < 0) {
      years--;
      months += 12;
    }
    let result = years + "г";
    if (months > 0) {
      result += " " + months + "м"
    }
    return result;
  };

  async presentDeleteConfirm() {
    const alert = await this.alertController.create({
      header: 'Потвърди изтриване',
      message: 'Сигурни ли сте, че искате да изтриете този запис?',
      buttons: [
        {
          text: 'Отказ',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {}
        },
        {
          text: 'Изтрий',
          handler: () => {
            this.patientService.deletePatient(this.patient!.id).subscribe({
              next: _ => {
                this.patient = undefined;
                this.goBack();
              }
            })
          }
        }
      ]
    });

    await alert.present();
  }

}
