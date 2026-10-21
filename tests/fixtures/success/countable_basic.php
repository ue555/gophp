<?php
class CountingBox implements Countable {
    private int $calls = 0;
    public function count(): int {
        $this->calls++;
        return $this->calls;
    }
}
$box = new CountingBox;
echo count($box), '|', count($box, COUNT_RECURSIVE), '|', $box->count(), '|';
class ChildCountingBox extends CountingBox {}
echo count(new ChildCountingBox);
