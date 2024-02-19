import { Injectable } from '@angular/core';
import { NodeInfoModule } from './node-info.module';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { nodeInfoRsp } from '../../sharable';

@Injectable({
  providedIn: NodeInfoModule
})
export class NodeInfoClientService {

  private  uri:string = "/api/v1/NodeDetails"
  constructor(private http:HttpClient) { }
  getInfo():Observable<nodeInfoRsp> {
    return  this.http.get(this.uri) as  Observable<nodeInfoRsp>
  }
}
