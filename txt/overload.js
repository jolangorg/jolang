class V {

    addLocal(){
        if (arguments.length === 1 && arguments[0] instanceof Vec2) return this.addLocal$Vec2(...arguments);
        if (arguments.length === 2 && typeof arguments[0] === "number" && typeof arguments[1] === "number") return this.addLocal$floatfloat(...arguments);

        return super.addLocal(...arguments);
    }

    addLocal$Vec2(v) {
        this.x = v.x;
        this.y = v.y;
        return this;
    }

    addLocal$floatfloat(x, y) {
        this.x = this.x;
        this.y = this.y;
        return this;
    }

}