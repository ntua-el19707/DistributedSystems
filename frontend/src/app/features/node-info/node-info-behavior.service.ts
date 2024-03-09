import { Injectable } from '@angular/core';
import { NodeInfoModule } from './node-info.module';
import { NodeInfoClientService } from './node-info-client.service';
import { NodeInfoSelfGraphQL, nodeDetails, nodeInfoRsp } from '../../sharable';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: NodeInfoModule,
})
export class NodeInfoBehaviorService {
  #nodeInfoBehaviorSubject: BehaviorSubject<NodeInfoSelfGraphQL> =
    new BehaviorSubject<NodeInfoSelfGraphQL>({
      client: { indexId: -1, nodeId: '', uri: '', uriPublic: '' },
      total: 0,
    });

  constructor(private nodeInfoClientService: NodeInfoClientService) {}
  fetchNodeInfo() {
    this.nodeInfoClientService.getInfo().subscribe(
      (r) => {
        const node = r.data?.self;
        if (node) {
          this.#nodeInfoBehaviorSubject.next(node);
        }
      },
      (err) => {},
      () => {}
    );
  }
  getNodeInfo(): BehaviorSubject<NodeInfoSelfGraphQL> {
    return this.#nodeInfoBehaviorSubject;
  }
}
