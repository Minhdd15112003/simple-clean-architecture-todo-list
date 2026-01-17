import { Handle, type Node, NodeProps, Position, useReactFlow } from "@xyflow/react";
import { ChangeEvent, useCallback, useState } from "react";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import { INodedata } from "@/types/flow";
import { Input } from "../ui/input";
import { Switch } from "../ui/switch";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Slider } from "../ui/slider";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Button } from "../ui/button";

const DefaultNode = ({
  children,
  data,
  style,
}: {
  children: React.ReactNode;
  data: any;
  style?: string;
}) => {
  return (
    <Card className={`w-48 min-h-24 border-2 border-gray-300 ${style ? style : ""}`}>
      <Handle type='target' position={Position.Top} id='top' className='w-3 h-3 bg-blue-500' />
      <Handle
        type='source'
        position={Position.Bottom}
        id='bottom'
        className='w-3 h-3 bg-blue-500'
      />
      <Handle type='target' position={Position.Left} id='left' className='w-3 h-3 bg-blue-500' />
      <Handle type='source' position={Position.Right} id='right' className='w-3 h-3 bg-blue-500' />
      <CardHeader>
        <CardTitle>
          {data.label.replace(data.label.at(0), data.label.at(0).toUpperCase())}
        </CardTitle>
      </CardHeader>
      <CardContent>{children}</CardContent>
    </Card>
  );
};
const TIME_OUT = 1000;
const nodeTypes = {
  start: function start(props: NodeProps) {
    const [isStarted, setIsStarted] = useState(false);

    const handleToggle = () => {
      const newState = !isStarted;
      setIsStarted(newState);
      console.log(newState); // "Returns" the new state via console log
    };

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-4'>
          <Button
            onClick={handleToggle}
            variant={isStarted ? "destructive" : "default"}
            className={`w-full text-white ${!isStarted ? "bg-green-500 hover:bg-green-600" : ""}`}
          >
            {isStarted ? "Stop" : "Start"}
          </Button>
          <p className='text-sm text-muted-foreground'>Điểm bắt đầu của flow</p>
        </div>
      </DefaultNode>
    );
  },
  // Screenshot node - không cần input
  screenshot: function screenshot(props: NodeProps) {
    return (
      <DefaultNode data={props.data}>
        <p className='text-sm text-muted-foreground'>Chụp màn hình thiết bị</p>
      </DefaultNode>
    );
  },

  // Sleep node
  sleep: function sleep(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Time sleep</Label>
            <Input
              type='number'
              placeholder='Thời gian ngủ (ms)...'
              onChange={(e) => updateNodeData(props.id, { timeSleep: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              placeholder='Nhập thời gian...'
              defaultValue={TIME_OUT}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Toast node
  toast: function toast(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div>
          <Label>Text</Label>
          <Input
            type='text'
            placeholder='Nội dung toast...'
            onChange={(e) => updateNodeData(props.id, { text: e.target.value })}
          />
        </div>
        <div>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Unlock Screen
  unlockScreen: function unlockScreen(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            placeholder='Timeout...'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Press Home
  pressHome: function pressHome(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            placeholder='Timeout...'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Press Back
  pressBack: function pressBack(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            placeholder='Timeout...'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Press Switch
  pressSwitch: function pressSwitch(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            placeholder='Timeout...'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Get Clipboard
  getClipboard: function getClipboard(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Timeout (ms)</Label>
          <Input
            type='number'
            placeholder='Timeout...'
            defaultValue={TIME_OUT}
            onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Set Clipboard
  setClipboard: function setClipboard(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Text</Label>
            <Textarea
              placeholder='Nhập text để copy...'
              className='min-h-[80px]'
              onChange={(e) => updateNodeData(props.id, { text: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={TIME_OUT}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // App Start
  appStart: function appStart(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Package Name</Label>
            <Input
              type='text'
              placeholder='com.example.app'
              onChange={(e) => updateNodeData(props.id, { packageName: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={5000}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // App Stop
  appStop: function appStop(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Package Name</Label>
            <Input
              type='text'
              placeholder='com.example.app'
              onChange={(e) => updateNodeData(props.id, { packageName: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={3000}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // App Clear
  appClear: function appClear(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Package Name</Label>
            <Input
              type='text'
              placeholder='com.example.app'
              onChange={(e) => updateNodeData(props.id, { packageName: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={3000}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Swipe
  swipe: function swipe(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div className='grid grid-cols-2 gap-2'>
            <div>
              <Label>Start X</Label>
              <Input
                type='number'
                placeholder='Start X'
                onChange={(e) => updateNodeData(props.id, { startX: e.target.value })}
              />
            </div>
            <div>
              <Label>Start Y</Label>
              <Input
                type='number'
                placeholder='Start Y'
                onChange={(e) => updateNodeData(props.id, { startY: e.target.value })}
              />
            </div>
          </div>
          <div className='grid grid-cols-2 gap-2'>
            <div>
              <Label>End X</Label>
              <Input
                type='number'
                placeholder='End X'
                onChange={(e) => updateNodeData(props.id, { endX: e.target.value })}
              />
            </div>
            <div>
              <Label>End Y</Label>
              <Input
                type='number'
                placeholder='End Y'
                onChange={(e) => updateNodeData(props.id, { endY: e.target.value })}
              />
            </div>
          </div>
          <div>
            <Label>Duration (ms)</Label>
            <Input
              type='number'
              onChange={(e) => updateNodeData(props.id, { duration: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={TIME_OUT}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Key Code
  keyCode: function keyCode(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Key Code</Label>
            <Input
              type='number'
              placeholder='Android key code'
              onChange={(e) => updateNodeData(props.id, { keyCode: e.target.value })}
            />
          </div>
          <div>
            <Label>Meta State</Label>
            <Input
              type='number'
              placeholder='Meta state (optional)'
              onChange={(e) => updateNodeData(props.id, { metaState: e.target.value })}
            />
          </div>
          <div className='flex items-center space-x-2'>
            <Switch
              id='repeat'
              onCheckedChange={(checked) => updateNodeData(props.id, { repeat: checked })}
            />
            <Label htmlFor='repeat'>Repeat</Label>
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={TIME_OUT}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Key Text
  keyText: function keyText(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Text</Label>
            <Input
              type='text'
              placeholder='Nhập text...'
              onChange={(e) => updateNodeData(props.id, { text: e.target.value })}
            />
          </div>
          <div>
            <Label>Timeout (ms)</Label>
            <Input
              type='number'
              defaultValue={TIME_OUT}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Device ID
  deviceId: function deviceId(props: NodeProps) {
    return (
      <DefaultNode data={props.data}>
        <p className='text-sm text-muted-foreground'>Lấy Device ID</p>
      </DefaultNode>
    );
  },

  // Resolution
  resolution: function resolution(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <div>
            <Label>Width</Label>
            <Input
              type='number'
              placeholder='width'
              defaultValue={0}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
          <div>
            <Label>Height</Label>
            <Input
              type='number'
              placeholder='height'
              defaultValue={0}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
          <div>
            <Label>orientation</Label>
            <Input
              type='number'
              defaultValue={720}
              onChange={(e) => updateNodeData(props.id, { timeout: e.target.value })}
            />
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Find Elements
  findElements: function findElements(props: NodeProps) {
    const { updateNodeData } = useReactFlow();
    const [action, setAction] = useState<string>("click");
    const [strategy, setStrategy] = useState<string>("resourceId");

    const handleStrategyChange = (value: string) => {
      setStrategy(value);
      // Không lưu strategy, chỉ lưu key là strategy name
    };

    const handleActionChange = (value: string) => {
      setAction(value);
      // Không lưu action, sẽ lưu key là action name với params
    };

    const handleValueChange = (value: string) => {
      // Lưu value với key là strategy hiện tại (resourceId, text, xpath, className)
      updateNodeData(props.id, { [strategy]: value });
    };

    const updateActionParams = (actionType: string, params: any) => {
      // Lưu params với key là action name (click, setText, all)
      updateNodeData(props.id, { [actionType]: params });
    };

    return (
      <DefaultNode data={props.data} style='w-80'>
        <div className='grid grid-cols-2 grid-rows-1 gap-4'>
          <div className='Strategy space-y-2'>
            <Label>Strategy</Label>
            <Select value={strategy} onValueChange={handleStrategyChange}>
              <SelectTrigger>
                <SelectValue placeholder='Chọn strategy' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='resourceId'>Resource ID</SelectItem>
                <SelectItem value='text'>Text</SelectItem>
                <SelectItem value='xpath'>XPath</SelectItem>
                <SelectItem value='className'>Class Name</SelectItem>
              </SelectContent>
            </Select>

            <Label>Value</Label>
            <Textarea
              placeholder='Nhập value...'
              className='min-h-[180px]'
              onChange={(e) => handleValueChange(e.target.value)}
            />
          </div>

          <div className='Action space-y-2'>
            <div className='space-y-2'>
              <Label>Action</Label>
              <Select value={action} onValueChange={handleActionChange}>
                <SelectTrigger>
                  <SelectValue placeholder='Chọn action' />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='click'>Click</SelectItem>
                  <SelectItem value='setText'>Set Text</SelectItem>
                  <SelectItem value='all'>Get All</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Parameters cho click */}
            {action === "click" && (
              <div className='space-y-2 bg-muted/50 p-2 rounded'>
                <div>
                  <Label className='text-xs'>Index (optional)</Label>
                  <Input
                    type='number'
                    placeholder='Element index'
                    defaultValue={0}
                    onChange={(e) =>
                      updateActionParams("click", {
                        ...(props.data.click || {}),
                        index: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Retry Time (optional)</Label>
                  <Input
                    type='number'
                    placeholder='Retry time'
                    defaultValue={3}
                    onChange={(e) =>
                      updateActionParams("click", {
                        ...(props.data.click || {}),
                        retry_time: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Timeout (ms)</Label>
                  <Input
                    type='number'
                    placeholder='Timeout'
                    defaultValue={TIME_OUT}
                    onChange={(e) =>
                      updateActionParams("click", {
                        ...(props.data.click || {}),
                        timeout: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
              </div>
            )}

            {/* Parameters cho setText */}
            {action === "setText" && (
              <div className='space-y-2 bg-muted/50 p-2 rounded'>
                <div>
                  <Label className='text-xs'>Text *</Label>
                  <Input
                    type='text'
                    placeholder='Nhập text...'
                    onChange={(e) =>
                      updateActionParams("setText", {
                        ...(props.data.setText || {}),
                        text: e.target.value,
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Index (optional)</Label>
                  <Input
                    type='number'
                    placeholder='Element index'
                    defaultValue={0}
                    onChange={(e) =>
                      updateActionParams("setText", {
                        ...(props.data.setText || {}),
                        index: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Retry Time (optional)</Label>
                  <Input
                    type='number'
                    placeholder='Retry time'
                    defaultValue={3}
                    onChange={(e) =>
                      updateActionParams("setText", {
                        ...(props.data.setText || {}),
                        retry_time: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Timeout (ms)</Label>
                  <Input
                    type='number'
                    placeholder='Timeout'
                    defaultValue={TIME_OUT}
                    onChange={(e) =>
                      updateActionParams("setText", {
                        ...(props.data.setText || {}),
                        timeout: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
              </div>
            )}

            {/* Parameters cho all */}
            {action === "all" && (
              <div className='space-y-2 bg-muted/50 p-2 rounded'>
                <div>
                  <Label className='text-xs'>Retry Time (optional)</Label>
                  <Input
                    type='number'
                    placeholder='Retry time'
                    defaultValue={3}
                    onChange={(e) =>
                      updateActionParams("all", {
                        ...(props.data.all || {}),
                        retry_time: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <div>
                  <Label className='text-xs'>Timeout (ms)</Label>
                  <Input
                    type='number'
                    placeholder='Timeout'
                    defaultValue={TIME_OUT}
                    onChange={(e) =>
                      updateActionParams("all", {
                        ...(props.data.all || {}),
                        timeout: parseInt(e.target.value),
                      })
                    }
                  />
                </div>
                <p className='text-xs text-muted-foreground'>Trả về: Promise&lt;any[]&gt;</p>
              </div>
            )}
          </div>
        </div>
      </DefaultNode>
    );
  },

  // Loop node (giữ lại từ code cũ)
  loop: function loop(props: NodeProps) {
    const { updateNodeData } = useReactFlow();

    return (
      <DefaultNode data={props.data}>
        <div className='space-y-2'>
          <Label>Iterations</Label>
          <Input
            type='number'
            placeholder='Số lần lặp...'
            defaultValue={1}
            onChange={(e) => updateNodeData(props.id, { iterations: e.target.value })}
          />
        </div>
      </DefaultNode>
    );
  },

  // Pause node (giữ lại từ code cũ)
  pause: function pause(props: NodeProps) {
    const [isStarted, setIsStarted] = useState(false);
    const { updateNodeData } = useReactFlow();

    const handleToggle = () => {
      const newState = !isStarted;
      setIsStarted(newState);
      updateNodeData(props.id, { isPaused: newState });
    };
    return (
      <DefaultNode data={props.data}>
        <div className='flex items-center space-x-2'>
          <Button
            onClick={handleToggle}
            variant={isStarted ? "destructive" : "default"}
            className={`w-full text-white ${!isStarted ? "bg-green-500 hover:bg-green-600" : ""}`}
          >
            {isStarted ? "Stop" : "Start"}
          </Button>
        </div>
      </DefaultNode>
    );
  },

  // setText: function setText(props: NodeProps) {
  //   return (
  //     <DefaultNode data={props.data}>
  //       <div className='space-y-2'>
  //         <div>
  //           <Label>Text</Label>
  //           <Textarea
  //             placeholder='Nhập text...'
  //             className='min-h-[80px]'
  //             onChange={(e) => console.log(e.target.value)}
  //           />
  //         </div>
  //         <div>
  //           <Label>Timeout (ms)</Label>
  //           <Input type='number' defaultValue={TIME_OUT} />
  //         </div>
  //       </div>
  //     </DefaultNode>
  //   );
  // },

  // click: function click(props: NodeProps) {
  //   return (
  //     <DefaultNode data={props.data}>
  //       <div className='space-y-2'>
  //         <div className='grid grid-cols-2 gap-2'>
  //           <div>
  //             <Label>X</Label>
  //             <Input type='number' placeholder='X coordinate' />
  //           </div>
  //           <div>
  //             <Label>Y</Label>
  //             <Input type='number' placeholder='Y coordinate' />
  //           </div>
  //         </div>
  //         <div>
  //           <Label>Duration (ms)</Label>
  //           <Input type='number' defaultValue={100} />
  //         </div>
  //         <div>
  //           <Label>Timeout (ms)</Label>
  //           <Input type='number' defaultValue={TIME_OUT} />
  //         </div>
  //       </div>
  //     </DefaultNode>
  //   );
  // },
};
export default nodeTypes;

// declare interface ElementWrapper {
//   click({
//     index,
//     retry_time,
//     timeout,
//   }: {
//     index?: number;
//     retry_time?: number;
//     timeout?: number;
//   }): Promise<void>;
//   setText({
//     text,
//     index,
//     retry_time,
//     timeout,
//   }: {
//     text: string;
//     index?: number;
//     retry_time?: number;
//     timeout?: number;
//   }): Promise<void>;
//   all({ retry_time, timeout }: { retry_time?: number; timeout?: number }): Promise<any[]>;
// }
// declare interface DeviceController {
//   loop(iterations: number, callback: (index: number) => Promise<void>): Promise<void>;
//   pause(is_pause: boolean): void;
//   deviceId(): string;
//   resolution(): { width: number; height: number; orientation: number };
//   screenshot(): string;
//   sleep(timeout?: number): Promise<any>;
//   toast(text: string, timeout?: number): Promise<any>;
//   unlockScreen(timeout?: number): Promise<any>;
//   pressHome(timeout?: number): Promise<any>;
//   pressBack(timeout?: number): Promise<any>;
//   pressSwitch(timeout?: number): Promise<any>;
//   getClipboard(timeout?: number): Promise<string>;
//   setClipboard(text: string, timeout?: number): Promise<any>;
//   appStart(packageName: string, timeout?: number): Promise<any>;
//   appStop(packageName: string, timeout?: number): Promise<any>;
//   appClear(packageName: string, timeout?: number): Promise<any>;
//   click(x: number, y: number, duration?: number, timeout?: number): Promise<any>;
//   swipe(
//     start_x: number,
//     start_y: number,
//     end_x: number,
//     end_y: number,
//     duration: number,
//     timeout?: number,
//   ): Promise<any>;
//   keyCode(key: number, meta_state?: number, repeat?: boolean, timeout?: number): Promise<any>;
//   keyText(text: string, timeout?: number): Promise<any>;
//   setText(text: string, timeout?: number): Promise<any>;
//   findElements(
//     strategy: "resourceId" | "text" | "xpath" | "className",
//     value: string,
//   ): ElementWrapper;
//   dumpJson(timeout?: number): Promise<any>;
// }
