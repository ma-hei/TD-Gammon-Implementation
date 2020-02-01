import { BoardDimensions } from './board-dimensions';

export class Point {
    public pointsString: string;
    public fillString: string;

    constructor(
        public id: number,
        boardDimensions: BoardDimensions
    ) {
      this.pointsString = BoardDimensions.computePointPositionString(id, boardDimensions);
      this.fillString = BoardDimensions.getPointColor(id);
    }
 } 
