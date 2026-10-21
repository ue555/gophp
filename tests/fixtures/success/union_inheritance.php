<?php
interface UnionContract { public function choose(int|string $value): int|string; }
class UnionParent implements UnionContract {
    public function choose(int|string $value): int|string { return $value; }
}
class UnionChild extends UnionParent {
    public function choose(string|int $value): string|int { return $value; }
}
echo gettype((new UnionChild)->choose(5)), '|', gettype((new UnionChild)->choose('a')), "\n";
