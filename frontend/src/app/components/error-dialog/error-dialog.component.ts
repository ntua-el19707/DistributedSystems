import { Component } from '@angular/core';
import { Inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import {

  MAT_DIALOG_DATA, MatDialogModule, MatDialogRef,

} from '@angular/material/dialog';
@Component({
  selector: 'app-error-dialog',
  standalone: true,
  imports: [MatDialogModule ,MatCardModule,MatButtonModule],
  templateUrl: './error-dialog.component.html',
  styleUrl: './error-dialog.component.scss'
})
export class ErrorDialogComponent {
errMsg:string = ""
constructor(
    @Inject(MAT_DIALOG_DATA) private data: string,
    private dialogRef: MatDialogRef<ErrorDialogComponent>){
      this.errMsg = data
    }
      closeDialog(): void {
    this.dialogRef.close();
  }
}
