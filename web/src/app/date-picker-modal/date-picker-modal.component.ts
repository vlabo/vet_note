import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'datepicker',
  templateUrl: './date-picker-modal.component.html',
  styleUrls: ['./date-picker-modal.component.scss'],
})
export class DatePickerModalComponent implements OnInit {

  @Input() value!: string;
  @Input() disabled: boolean = false;
  @Output() valueChange = new EventEmitter<string>();

  isPopoverOpen = false;
  selectedDate: string = "";

  ngOnInit(): void {
    this.selectedDate = this.value;
  }

  openDatePickerPopover() {
    this.selectedDate = this.value;
    this.isPopoverOpen = true;
  }

  onPopoverDismiss() {
    this.isPopoverOpen = false;
  }

  onDateChange(event: any) {
    this.selectedDate = event.detail.value;
  }

  confirmDate() {
    this.value = this.selectedDate;
    this.valueChange.emit(this.value);
    this.onPopoverDismiss();
  }

  isValidDate(dateString: string): boolean {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  }
}
