import { Component, OnInit } from '@angular/core';
import { PatientsService, Procedure } from '../patients.service';
import { ActivatedRoute, Router } from '@angular/router';
import { IonicModule } from '@ionic/angular';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-edit-procedure',
  templateUrl: './edit-procedure.component.html',
  styleUrls: ['./edit-procedure.component.scss'],
  imports: [IonicModule, CommonModule, FormsModule],
  standalone: true,
})
export class EditProcedureComponent implements OnInit {
  procedure: Procedure = new Procedure(); //{ Type: '', Date: '', Details: '' };
  isEditMode: boolean = false;
  procedureTypes: string[] = ['Examination', 'Surgery', 'Vaccine', 'Castration', 'Blood Test']; // Predefined procedure types
  patientId: string | null = null;
  date: string = "";

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private patientsService: PatientsService
  ) { }

  ngOnInit(): void {
    console.log("Patient Edit init");
    const procedureId = this.route.snapshot.paramMap.get('procedureId');
    if (!procedureId) {
      return;
    }

    let procedure = this.patientsService.getProcedure(this.patientId!);
    if(!procedure) {
      return;
    }

    this.procedure = { ...procedure};
    this.date = this.procedure.Date.toISOString();
    this.isEditMode = true;
  }

  goBack(): void {
    if (this.patientId) {
      this.router.navigate(['/patient', this.patientId]);
    }
  }

  saveProcedure(): void {
    // TODO: Update data.
    this.patientsService.updateProcedure(this.procedure);
  }
}
