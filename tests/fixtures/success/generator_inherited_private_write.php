<?php
class ParentOutput {
    private int $verbosity;
    public function setVerbosity(int $value): void {
        $this->verbosity = $value;
    }
    public function readVerbosity(): int {
        return $this->verbosity;
    }
}
class ChildOutput extends ParentOutput {}
$output = new ChildOutput;
$output->setVerbosity(3);
echo $output->readVerbosity();
