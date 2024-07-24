import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { IonicModule } from '@ionic/angular';
import { CommonModule, Location } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PatientsService } from '../patients.service';

import { addIcons } from "ionicons";
import { arrowBack } from "ionicons/icons";
import { Procedure } from '../../../server/bindings/Procedure';

@Component({
  selector: 'app-procedure',
  templateUrl: './procedure.component.html',
  styleUrls: ['./procedure.component.scss'],
  imports: [IonicModule, CommonModule, FormsModule],
  standalone: true,
})
export class ProcedureComponent implements OnInit {
  procedure: Procedure | null = null;
  isEditMode: boolean = false;
  isNewMode: boolean = false;

  procedureTypes: string[] = ['Examination', 'Surgery', 'Vaccine', 'Castration', 'Blood Test'];
  patientId: string | null = null;
  date: string = "";

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private location: Location,
    private patientsService: PatientsService
  ) { 
    addIcons({"arrow-back": arrowBack})
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(paramMap => {
      const procedureId = paramMap.get('procedureId');
      const patientId = paramMap.get('patientId');

      if (patientId) {
        this.patientId = patientId;
        this.isNewMode = true;
        this.isEditMode = true;
      } else {
        this.isEditMode = false;
      }
      if (procedureId) {
        const procedure = this.patientsService.getProcedure(procedureId);
          if (procedure) {
            this.procedure = { ...procedure };
            console.log(this.procedure);
            this.date = this.procedure.date;
          }
        }
    });

    this.route.queryParamMap.subscribe(queryParams => {
      if(!this.isNewMode) {
        this.isEditMode = queryParams.get('edit') === 'true';
      }
    });
  }

  goBack(): void {
    this.location.back();
  }

  enableEditMode() {
    this.isEditMode = true;
    this.router.navigate([], { queryParams: { edit: true }, queryParamsHandling: 'merge' });
  }

  saveProcedure(): void {
    // TODO: update date
    if (this.isNewMode) {
      this.patientsService.addProcedure(this.patientId!, this.procedure!);
      this.router.navigate(["procedure", this.procedure!.id], { queryParams: { edit: false },  replaceUrl: true });
    } else if (this.isEditMode) {
      this.patientsService.updateProcedure(this.procedure!);
      this.goBack();
    }
  }
}