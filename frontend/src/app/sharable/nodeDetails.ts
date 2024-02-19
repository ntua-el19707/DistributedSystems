export interface  nodeDetails  {
    nodeId:  string 
    indexId:number 
    Uri:string
}


export interface  nodeInfoRsp {

    client: nodeDetails
    total:number
} 