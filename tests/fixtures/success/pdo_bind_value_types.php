<?php
$db = new PDO('sqlite::memory:');
foreach ([
    [123, PDO::PARAM_STR],
    ['42', PDO::PARAM_INT],
    [true, PDO::PARAM_BOOL],
    [7, PDO::PARAM_BOOL],
    ['ignored', PDO::PARAM_NULL],
    ["A\0B", PDO::PARAM_LOB],
    ['', PDO::PARAM_LOB],
] as $case) {
    $stmt = $db->prepare('SELECT typeof(:data) AS kind, hex(:data) AS bytes');
    $stmt->bindValue(':data', $case[0], $case[1]);
    $stmt->execute();
    $row = $stmt->fetch(PDO::FETCH_ASSOC);
    echo $row['kind'], ':', $row['bytes'], '|';
}
$stmt = $db->prepare('SELECT typeof(?)');
$stmt->bindValue(1, 321);
$stmt->execute();
echo $stmt->fetchColumn(), '|';
$stmt->execute([456]);
echo $stmt->fetchColumn(), '|';
$stmt->execute();
echo $stmt->fetchColumn(), '|';
$stream = fopen('php://temp', 'w+');
fwrite($stream, "C\0D");
fseek($stream, 0);
$stmt = $db->prepare('SELECT typeof(?) AS kind, hex(?) AS bytes');
$stmt->bindValue(1, $stream, PDO::PARAM_LOB);
$stmt->bindValue(2, "C\0D", PDO::PARAM_LOB);
$stmt->execute();
$row = $stmt->fetch(PDO::FETCH_ASSOC);
echo $row['kind'], ':', $row['bytes'];
fclose($stream);
class BoundText {
    public int $calls = 0;
    public function __toString(): string {
        $this->calls++;
        return 'ready';
    }
}
$text = new BoundText;
$stmt = $db->prepare('SELECT hex(?)');
$stmt->bindValue(1, $text);
echo '|', $text->calls;
$stmt->execute();
echo '|', $text->calls, '|', $stmt->fetchColumn();
$stmt = $db->prepare('SELECT typeof(:named)');
$stmt->bindValue(type: PDO::PARAM_INT, value: '87', param: ':named');
$stmt->execute();
echo '|', $stmt->fetchColumn();
