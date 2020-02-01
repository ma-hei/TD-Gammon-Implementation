export class BoardDimensions {

  boardWidth: number;
  boardHeight: number;
  leftOffset: number = 10;
  pointHeight: number = 200;
  pointSpacing: number;
  pointWidth: number;
  a: number = 0.2; //pointSpacing = a*pointWidth

  constructor(
    boardWidth: number,
    boardHeight: number
  ) 
  {
    this.boardWidth = boardWidth;
    this.boardHeight = boardHeight;
    this.pointSpacing = this.computePointSpacing();
    this.pointWidth = this.computePointWidth();
  }

  // width - 2*offset = 13*pointW + 12*(a*pointW);
  //                  = pointW * (13 + 12a);
  //                  = 13*(pointS/a) + 12*pointS
  //                  = pointS * (12 + 13/a)
  computePointSpacing(): number {
    return (this.boardWidth - 2 * this.leftOffset) / (12 + (13/this.a));
  }   

  computePointWidth(): number {
    return (this.boardWidth - 2 * this.leftOffset) / (13 + 12 * this.a);
  }

  static generateSinglePointString(x: number, y: number) : string {
    return x.toString() + "," + y.toString() + " "; 
  }

  static computePointPositionString(pointId: number, boardDimensions: BoardDimensions): string {
    let lowerY;
    let upperY;
    if (pointId > 11) {
      lowerY = 0;
      upperY = 0 + boardDimensions.pointHeight;
    } else {
      lowerY = boardDimensions.boardHeight;
      upperY = lowerY - boardDimensions.pointHeight;
    }
    let startX = boardDimensions.leftOffset + (pointId % 12) * (boardDimensions.pointSpacing + boardDimensions.pointWidth);
    if ((pointId > 5 && pointId < 12) || (pointId > 17)) {
      startX += boardDimensions.pointWidth;
      startX += boardDimensions.pointSpacing;
    }
    let midpointX = startX + (0.5 * boardDimensions.pointWidth);
    let endpointX = startX + boardDimensions.pointWidth;
    let temp = this.generateSinglePointString(startX, lowerY) + 
               this.generateSinglePointString(midpointX, upperY) + 
               this.generateSinglePointString(endpointX, lowerY);
    return temp; 
  }

  static getPointColor(id: number) : string {
      let color;
      if (id % 2 == 0) {
        color = "#800000";
      } else
      {
        color = "#666699";
      }
      return color;
    }

  getBarX() : number {
    return this.leftOffset + 6*(this.pointWidth + this.pointSpacing) + 0.5*this.pointWidth;
  }
  
  static computeCheckerX(pointId: number, boardDimensions: BoardDimensions) : number {
    let x = boardDimensions.leftOffset + (pointId % 12) * (boardDimensions.pointSpacing + boardDimensions.pointWidth) + 0.5 * boardDimensions.pointWidth;
    if ((pointId > 5 && pointId < 12 ) || pointId > 17) {
      x+=boardDimensions.pointWidth + boardDimensions.pointSpacing;
    }
    return x;
  }

  static computeCheckerY(pointId: number, nCheckersOnPoint: number, boardDimensions: BoardDimensions) {
    let radius = boardDimensions.pointHeight / 9;
    let offset = (nCheckersOnPoint * (radius * 2) + radius);
    let y;
    if (pointId > 11) {
      y = offset; 
    } else
    {
      y = boardDimensions.boardHeight - offset;
    }
    return y;
  }
  
  static computeCheckerRadius(boardDimensions: BoardDimensions) : number {
    return boardDimensions.pointHeight / 9;
  }

  static getCheckerColor(team: number) : string {
    if (team==1) {
      return "#ffffff";
    }
    else {
      return "#ff3300";
    }
  }

}
