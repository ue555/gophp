<?php
class ChildFailure extends Exception {
    public function __construct($previous) {
        parent::__construct('child', 7, $previous);
    }
}
$first = new Error('root');
$child = new ChildFailure($first);
echo $child->getMessage(), '|', $child->getCode(), '|';
echo ($child->getPrevious() === $first ? 'same' : 'wrong'), "\n";
