<?php
class UnionParent { public function take(int|string $value): int|string { return $value; } }
class UnionChild extends UnionParent { public function take(array $value): array { return $value; } }
