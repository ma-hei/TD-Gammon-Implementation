export class Point {
    public pointsString: string;
    spacing: number;
    width: number;
    height: number;
    lowerY: number;
    upperY: number;
    leftOffset: number;
    
    constructor(
        public id: number,
        public boardWidth: number,
        public boardHeight: number 
    ) {
      // 11 * x + y
      // 11 * (a+y) + y
      // 12y + 11 a = width -> 23y
      this.leftOffset = 10;
      this.spacing = boardWidth / 23;
      this.width = this.spacing;
      this.lowerY = boardHeight;

      this.height = 50;
      this.upperY = this.lowerY - this.height;
      if (id > 11) {
        this.mirrorOnYAxis();
      }
      this.setPointsString(id);
    }
  
    generatePointString(x: number, y: number) : string {
      return x.toString() + "," + y.toString() + " "; 
    }

    mirrorOnYAxis() : void {
      this.lowerY = 0;
      this.upperY = 0 + this.height;
    }

    setPointsString(id: number) : void {
      let startX = this.leftOffset + (id % 12) * (this.spacing + this.width);
      let midpointX = startX + (0.5 * this.width);
      let endpointX = startX + this.width;
      let temp = this.generatePointString(startX, this.lowerY) + this.generatePointString(midpointX, this.upperY) + this.generatePointString(endpointX, this.lowerY);
      this.pointsString = temp;
    } 
}
