import { Injectable } from '@angular/core';
import * as uuid from 'uuid';


export class Procedure {
  Id: string = "";
  Type: string = "";
  Date: string = "";
  Details: String = "";
}

export class Patient {
  Id: string = "";

  Type: string = "";
  Name: string = "";
  Gender: "male" | "female" | "unknown" = "unknown";
  BirthDate: string = "";
  ChipId: string = "";
  Weight: number = 0;
  Castrated: boolean = false;
  LastModified: Date = new Date();

  Procedures: Procedure[] = [];
  Owner: string = "";
  OwnerPhone: string = "";
}

export class ListPatient {
  Id: string = "";
  Type: string = "";
  Name: string = "";
  ChipId: string = "";
  Owner: string = "";
  Phone: string = "";
}

export class ViewPatient {
  Id: string = "";

  Type: string = "";
  Name: string = "";
  Gender: "male" | "female" | "unknown" = "unknown";
  BirthDate: string = "";
  ChipId: string = "";
  Weight: number = 0;
  Castrated: boolean = false;

  Procedures: Procedure[] = [];
  Owner: string = "";
  OwnerPhone: string = "";
}


@Injectable({
  providedIn: 'root'
})
export class PatientsService {
  private patients = new Map<string, Patient>()
  private procedures = new Map<string, Procedure>()
  private types: String[] = ["Куче", "Котка", "Прица", "Заек"];

  private generateMockData() {
    this.procedures.set("0510ee92-b30e-4bff-a6d9-6af70b0e6acc", {
      Id: "0510ee92-b30e-4bff-a6d9-6af70b0e6acc",
      Type: 'Преглед',
      Date: '2023-01-15',
      Details: 'Routine check-up. All vitals are normal.'
    });
    this.procedures.set("b7306c26-8b48-4a8b-83c9-a2425c117364", {
      Id: "b7306c26-8b48-4a8b-83c9-a2425c117364",
      Type: 'Ваксина',
      Date: '2023-02-20',
      Details: 'Administered rabies vaccine.'
    });
    this.procedures.set("661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1", {
      Id: "661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1",
      Type: 'Кръвно изледване',
      Date: '2023-03-10',
      Details: 'Blood test for heartworm. Results are negative.'
    });

    let procedures = Array.from(this.procedures.values());
    let patient: Patient = {
      Id: "1",
      Type: "Куче",
      Name: "Buddy",
      Owner: "John Doe",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "male",
      BirthDate: "2018-01-15",
      ChipId: "482736194",
      LastModified: new Date("2023-01-01"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "2",
      Type: "Котка",
      Name: "Whiskers",
      Owner: "Jane Smith",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "female",
      BirthDate: "2019-05-20",
      ChipId: "193847562",
      LastModified: new Date("2023-02-15"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "3",
      Type: "Прица",
      Name: "Tweety",
      Owner: "Alice Johnson",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "unknown",
      BirthDate: "2020-07-30",
      ChipId: "758392014",
      LastModified: new Date("2023-03-10"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "4",
      Type: "Заек",
      Name: "Thumper",
      Owner: "Bob Brown",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "female",
      BirthDate: "2021-11-05",
      ChipId: "100237000236519",
      LastModified: new Date("2023-04-20"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "5",
      Type: "Куче",
      Name: "Max",
      Owner: "Charlie Davis",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "male",
      BirthDate: "2017-03-10",
      ChipId: "100237000236511",
      LastModified: new Date("2023-05-25"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "6",
      Type: "Куче",
      Name: "Бъди",
      Owner: "Иван Иванов",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "male",
      BirthDate: "2018-01-15",
      ChipId: "100237000236514",
      LastModified: new Date("2023-01-01"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "7",
      Type: "Котка",
      Name: "Мърка",
      Owner: "Мария Петрова",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "female",
      BirthDate: "2019-05-20",
      ChipId: "100237000236513",
      LastModified: new Date("2023-02-15"),
      Weight: 13,
      Castrated: false,
    };
    patient = {
      Id: "8",
      Type: "Прица",
      Name: "Чурулик",
      Owner: "Александър Георгиев",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "unknown",
      BirthDate: "2020-07-30",
      ChipId: "100237000236512",
      LastModified: new Date("2023-03-10"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "9",
      Type: "Rabbit",
      Name: "Тупър",
      Owner: "Борислав Димитров",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "female",
      BirthDate: "2021-11-05",
      ChipId: "100237000236517",
      LastModified: new Date("2023-04-20"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "10",
      Type: "Куче",
      Name: "Макс",
      Owner: "Георги Василев",
      OwnerPhone: "08773423453",
      Procedures: procedures,
      Gender: "male",
      BirthDate: "2017-03-10",
      ChipId: "100237000236513",
      LastModified: new Date("2023-05-25"),
      Weight: 13,
      Castrated: false,
    };
    this.patients.set(patient.Id, patient);
  }
  constructor() {
    this.generateMockData();
  }


  // public getPatients() {
  //   return Array.from(this.patients.values())
  // }

  public getPatient(key: string): ViewPatient | undefined {
    const patient = this.patients.get(key);
    return patient ? {
      Id: patient.Id,
      Type: patient.Type,
      Name: patient.Name,
      Gender: patient.Gender,
      BirthDate: patient.BirthDate,
      ChipId: patient.ChipId,
      Weight: patient.Weight,
      Castrated: patient.Castrated,

      Procedures: patient.Procedures,
      Owner: patient.Owner,
      OwnerPhone: patient.OwnerPhone,
    } : undefined;
  }

  public getProcedure(key: string): Procedure | undefined {
    const procedure = this.procedures.get(key);
    return procedure ? { ...procedure } : undefined;
  }

  public getProcedures(keys: string[]): Procedure[] {
    var procedures: Procedure[] = [];
    var service = this;
    keys.forEach(function(key) {
      const procedure = service.procedures.get(key);
      if (procedure) {
        procedures.push(procedure);
      }
    })

    return procedures;
  }

  public getPatientList(): ListPatient[] {
    var list: ListPatient[] = [];
    this.patients.forEach(function(p) {
      list.push({
        Id: p.Id,
        Type: p.Type,
        Name: p.Name,
        ChipId: p.ChipId,
        Owner: p.Owner,
        Phone: p.OwnerPhone,
      })
    });
    return list;
  }

  public addPatient(patient: ViewPatient) {
    patient.Id = uuid.v4();
    // this.patients.set(patient.Id, patient);
  }

  public addProcedure(patientId: string, procedure: Procedure) {
    procedure.Id = uuid.v4();
    this.procedures.set(procedure.Id, procedure);
    let patient = this.patients.get(patientId);
    patient?.Procedures.unshift(procedure);
  }

  public updatePatient(patient: ViewPatient) {
    // this.patients.set(patient.Id, patient);
  }

  public updateProcedure(procedure: Procedure) {
    this.procedures.set(procedure.Id, procedure);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
