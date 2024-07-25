import * as jo from "jo";

export class ShapeType extends jo.Enum {
    static CIRCLE = new ShapeType("CIRCLE");
    static EDGE = new ShapeType("EDGE");
    static POLYGON = new ShapeType("POLYGON");
    static CHAIN = new ShapeType("CHAIN");
}