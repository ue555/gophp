<?php
trait LeftMethod { public function run() { return 1; } }
trait RightMethod { public function run() { return 2; } }
echo "before\n";
class CollisionWork { use LeftMethod, RightMethod; }
echo "after\n";
