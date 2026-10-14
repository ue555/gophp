<?php
class LevelOne {
    private $value = 1;
    public function one() { return $this->value; }
}
class LevelTwo extends LevelOne {
    private $value = 2;
    public function two() { return $this->value; }
    public function changeTwo() { $this->value = 8; }
}
class LevelThree extends LevelTwo {
    private $value = 3;
    public function three() { return $this->value; }
}
$first = new LevelThree;
$second = new LevelThree;
echo $first->one(), '|', $first->two(), '|', $first->three(), '|';
echo $first == $second ? 'same' : 'different';
$first->changeTwo();
echo '|', $first->one(), '|', $first->two(), '|', $first->three(), '|';
echo $first == $second ? 'same' : 'different', "\n";
