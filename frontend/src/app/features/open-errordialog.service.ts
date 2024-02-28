import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ErrorDialogComponent } from '../components/error-dialog/error-dialog.component';

@Injectable({
  providedIn: 'root',
})
export class OpenErrordialogService {

  constructor(private dialog: MatDialog) { }
  errorDialog(errMsg :string ){

    this.dialog.open(ErrorDialogComponent , {data:errMsg,  disableClose: true,
    height:"300px"  , width:"400px"})
  }
}
