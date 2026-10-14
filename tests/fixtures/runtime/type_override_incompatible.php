<?php
class TypeParent { public function take($x): int { return 1; } }
class TypeChild extends TypeParent { public function take(int $x): int { return $x; } }
