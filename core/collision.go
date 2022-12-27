package core

import "platformer/common"

type Collider interface {
	GetCollisionBox() CollisionBox
}

type CollisionResult struct {
	newX       float64
	newY       float64
	hitFloor   bool
	hitCeiling bool
	hitWall    bool
}

func DoCollision(oldX, oldY, newX, newY, w, h float64, level *Level, tryFall bool) CollisionResult {

	cr := CollisionResult{
		newX: newX,
		newY: newY,
	}

	isFalling := newY > oldY

	horizontalTopLeftX, horizontalTopLeftY := newX, oldY
	horizontalTopRightX, horizontalTopRightY := newX+w, oldY
	horizontalBottomLeftX, horizontalBottomLeftY := newX, oldY+h
	horizontalBottomRightX, horizontalBottomRightY := newX+w, oldY+h
	verticalTopLeftX, verticalTopLeftY := oldX, newY
	verticalTopRightX, verticalTopRightY := oldX+w, newY
	verticalBottomLeftX, verticalBottomLeftY := oldX, newY+h
	verticalBottomRightX, verticalBottomRightY := oldX+w, newY+h

	var horizontalTopLeftCollider Collider = nil
	var horizontalTopRightCollider Collider = nil
	var horizontalBottomLeftCollider Collider = nil
	var horizontalBottomRightCollider Collider = nil
	var verticalTopLeftCollider Collider = nil
	var verticalTopRightCollider Collider = nil
	var verticalBottomLeftCollider Collider = nil
	var verticalBottomRightCollider Collider = nil
	for _, c := range level.GetColliders() {
		cb := c.GetCollisionBox()
		if common.Contains(cb.x, cb.y, cb.w, cb.h, horizontalTopLeftX, horizontalTopLeftY) {
			horizontalTopLeftCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, horizontalTopRightX, horizontalTopRightY) {
			horizontalTopRightCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, horizontalBottomLeftX, horizontalBottomLeftY) {
			horizontalBottomLeftCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, horizontalBottomRightX, horizontalBottomRightY) {
			horizontalBottomRightCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, verticalTopLeftX, verticalTopLeftY) {
			verticalTopLeftCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, verticalTopRightX, verticalTopRightY) {
			verticalTopRightCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, verticalBottomLeftX, verticalBottomLeftY) {
			verticalBottomLeftCollider = c
		}
		if common.Contains(cb.x, cb.y, cb.w, cb.h, verticalBottomRightX, verticalBottomRightY) {
			verticalBottomRightCollider = c
		}
	}

	// horizontal
	x, y := horizontalTopLeftX, horizontalTopLeftY
	td := level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block || horizontalTopLeftCollider != nil {
		cr.newX = oldX
		cr.hitWall = true
	}

	x, y = horizontalTopRightX, horizontalTopRightY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block || horizontalTopRightCollider != nil {
		cr.newX = oldX
		cr.hitWall = true
	}

	x, y = horizontalBottomLeftX, horizontalBottomLeftY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block || horizontalBottomLeftCollider != nil {
		cr.newX = oldX
		cr.hitWall = true
	}

	x, y = horizontalBottomRightX, horizontalBottomRightY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block || horizontalBottomRightCollider != nil {
		cr.newX = oldX
		cr.hitWall = true
	}

	// vertical
	x, y = verticalBottomLeftX, verticalBottomLeftY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block {
		cr.newY = float64(int(y/common.TileSize)*common.TileSize) - h - fudge
		cr.hitFloor = true
	}
	if verticalBottomLeftCollider != nil {
		cr.newY = verticalBottomLeftCollider.GetCollisionBox().y - h - fudge
		cr.hitFloor = true
	}
	if !tryFall && td.Platform && isFalling {
		distance := float64(int(y/common.TileSize)*common.TileSize) - (oldY + h)
		if distance > -1 {
			cr.newY = oldY + distance - fudge
			cr.hitFloor = true
		}
	}

	x, y = verticalTopLeftX, verticalTopLeftY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block {
		cr.newY = float64(int(y/common.TileSize)*common.TileSize) + common.TileSize + fudge
		if cr.newY > 0 {
			cr.hitCeiling = true
		}
	}
	if verticalTopLeftCollider != nil {
		cb := verticalTopLeftCollider.GetCollisionBox()
		cr.newY = cb.y + cb.h + fudge
		cr.hitFloor = true
	}

	x, y = verticalBottomRightX, verticalBottomRightY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block {
		cr.newY = float64(int(y/common.TileSize)*common.TileSize) - h - fudge
		cr.hitFloor = true
	}
	if verticalBottomRightCollider != nil {
		cr.newY = verticalBottomRightCollider.GetCollisionBox().y - h - fudge
		cr.hitFloor = true
	}
	if !tryFall && td.Platform && isFalling {
		distance := float64(int(y/common.TileSize)*common.TileSize) - (oldY + h)
		if distance > -1 {
			cr.newY = oldY + distance - fudge
			cr.hitFloor = true
		}
	}

	x, y = verticalTopRightX, verticalTopRightY
	td = level.tiledGrid.GetTileData(int(x/common.TileSize), int(y/common.TileSize))
	if td.Block {
		cr.newY = float64(int(y/common.TileSize)*common.TileSize) + common.TileSize + fudge
		if cr.newY > 0 {
			cr.hitCeiling = true
		}
	}
	if verticalTopRightCollider != nil {
		cb := verticalTopRightCollider.GetCollisionBox()
		cr.newY = cb.y + cb.h + fudge
		cr.hitFloor = true
	}

	return cr
}
