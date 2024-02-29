export interface  nodeDetails  {
    nodeId:  string 
    indexId:number 
    uri:string
    uriPublic:string
}


export interface  nodeInfoRsp {

    client: nodeDetails
    total:number
} 
export interface  clientsInfoRsp {

    clients: Array<nodeDetails>
    total:number
} 