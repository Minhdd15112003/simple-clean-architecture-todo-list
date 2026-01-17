import { INodedata, INodeposition } from "@/types/flow";
import { type Node } from "@xyflow/react";

class NewNode {
  readonly id: string;
  readonly position: INodeposition;
  readonly data: INodedata;
  readonly type: string;

  constructor(nodeLength: number, data: INodedata, position?: INodeposition, type?: string) {
    this.id = `${data.label}-${nodeLength + 1}`;
    if (position) {
      this.position = position;
    } else {
      this.position = { x: 0, y: 0 };
    }

    this.data = data;
    if (type) {
      this.type = type;
    } else {
      this.type = "default";
    }
  }
  /// viet cai nay thanh switch case
  CreateNode(): Node {
    return {
      id: this.id,
      position: this.position,
      data: this.data,
      type: this.data.label,
    };
  }
}
export default NewNode;
