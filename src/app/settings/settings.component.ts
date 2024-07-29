import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { IonicModule, ItemReorderEventDetail } from '@ionic/angular';
import { Location } from '@angular/common';

import { addIcons } from "ionicons";
import { arrowBack, add } from 'ionicons/icons';
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class SettingsComponent implements OnInit {

  constructor(private location: Location) { 
    addIcons({"add": add});
  }

  ngOnInit() {
    addIcons({
      "arrow-back": arrowBack
    })
  }

  handleReorder(ev: CustomEvent<ItemReorderEventDetail>) {
    // The `from` and `to` properties contain the index of the item
    // when the drag started and ended, respectively
    console.log('Dragged from index', ev.detail.from, 'to', ev.detail.to);

    // Finish the reorder and position the item in the DOM based on
    // where the gesture ended. This method can also be called directly
    // by the reorder group
    ev.detail.complete();
  }

  goBack() {
    this.location.back();
  }
}
