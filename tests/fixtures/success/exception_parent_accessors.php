<?php
class CustomError extends Exception {
    public function text() { return parent::getMessage(); }
    public function cause() { return parent::getPrevious(); }
    public function frames() { return parent::getTrace(); }
}
$first = new Error('root');
$error = new CustomError('hello', 0, $first);
echo $error->text(), '|', ($error->cause() === $first ? 'same' : 'wrong'), '|';
echo count($error->frames()), "\n";
