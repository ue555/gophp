<?php
class EvenFilter extends FilterIterator {
    public function accept(): bool {
        return $this->current() % 2 === 0;
    }
}
$iterator = new EvenFilter(new ArrayIterator([1, 2, 3, 4]));
foreach ($iterator as $key => $value) {
    echo $key, ':', $value, '|';
}
