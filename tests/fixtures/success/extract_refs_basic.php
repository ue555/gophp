<?php
$input = ['first' => 1, 'second' => 2];
$copy = $input;
echo extract(flags: EXTR_REFS, array: $input), '|';
$first = 9;
echo $input['first'], '|', $copy['first'], '|';
$input['second'] = 7;
echo $second, '|';
unset($first);
$input['first'] = 4;
echo $copy['first'], '|';

function localRefs($data) {
    extract($data, EXTR_REFS);
    $value = 6;
    return $data['value'];
}
echo localRefs(['value' => 1]), '|';

class RefCarrier {
    public $code = ['use' => ['value' => 3]];
    public function apply() {
        extract($this->code['use'], EXTR_REFS);
        $value = 8;
        return $this->code['use']['value'];
    }
}
echo (new RefCarrier)->apply(), '|';

$skip = ['kept' => 10, 'fresh' => 11];
$kept = 100;
echo extract($skip, EXTR_SKIP | EXTR_REFS), '|';
$fresh = 12;
echo $skip['fresh'], '|', $kept;
