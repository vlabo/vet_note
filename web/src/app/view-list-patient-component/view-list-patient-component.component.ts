import { ChangeDetectorRef, Component, ElementRef, OnInit, Renderer2 } from '@angular/core';

@Component({
  selector: 'app-view-list-patient-component',
  templateUrl: './view-list-patient-component.component.html',
  styleUrls: ['./view-list-patient-component.component.scss'],
})
export class ViewListPatientComponentComponent implements OnInit {

  constructor(
    private renderer: Renderer2,
    private el: ElementRef,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit() { }
  updateInnerHtml(newHtml: string) {
    this.renderer.setProperty(this.el.nativeElement, 'innerHTML', newHtml);
    this.cdr.detectChanges(); // Manually trigger change detection if needed
  }
}
