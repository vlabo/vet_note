use serde::{Deserialize, Serialize};
use ts_rs::TS;
use uuid::Uuid;

#[derive(Serialize, Deserialize, TS)]
#[ts(export)]
pub enum Gender {
    Male,
    Female,
    Unknown,
}

#[derive(Serialize, Deserialize, TS, Default)]
#[ts(export)]
pub struct Procedure {
    pub id: Uuid,
    pub r#type: String,
    pub date: String,
    pub details: String,
}

#[derive(Serialize, Deserialize, TS)]
#[ts(export)]
pub struct ListPatient {
    pub id: Uuid,
    pub r#type: String,
    pub name: String,
    pub chip_id: String,
    pub owner: String,
    pub phone: String,
}

#[derive(Serialize, Deserialize, TS)]
#[ts(export)]
pub struct Patient {
    pub id: Uuid,
    pub r#type: String,
    pub name: String,
    pub gender: Gender,
    pub birth_date: String,
    pub chip_id: String,
    pub weight: f64,
    pub castrated: bool,
    pub last_modified: String,
    pub note: String,
    pub owner: String,
    pub owner_phone: String,
    pub procedures: Vec<Procedure>,
}

impl Patient {
    fn as_view(&self) -> ListPatient {
        return ListPatient {
            id: self.id,
            r#type: self.r#type.clone(),
            name: self.name.clone(),
            chip_id: self.chip_id.clone(),
            owner: self.owner.clone(),
            phone: self.owner_phone.clone(),
        };
    }
}

impl Default for Patient {
    fn default() -> Self {
        Self {
            id: Uuid::default(),
            r#type: String::default(),
            name: String::default(),
            gender: Gender::Unknown,
            birth_date: String::default(),
            chip_id: String::default(),
            weight: 0.0,
            castrated: false,
            last_modified: String::default(),
            note: String::default(),
            owner: String::default(),
            owner_phone: String::default(),
            procedures: Vec::new(),
        }
    }
}

pub trait ToViewPatientList {
    fn as_view(&self) -> Vec<ListPatient>;
}

impl ToViewPatientList for Vec<Patient> {
    fn as_view(&self) -> Vec<ListPatient> {
        let mut list = Vec::new();
        for p in self {
            list.push(p.as_view());
        }
        return list;
    }
}
