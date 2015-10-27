#include <IRremote.h>

//set pins
#define REDPIN 9
#define GREENPIN 10
#define BLUEPIN 5

//global variables for colors
unsigned char Red = 0;
unsigned char Green = 0;
unsigned char Blue = 0;

//setting up the infrared receiver
int RECV_PIN = 2;
IRrecv irrecv(RECV_PIN);
decode_results results;

String inputString = "";         // a string to hold incoming data
boolean stringComplete = false;  // whether the string is complete

void setup()
{
  Serial.begin(9600);
  attachInterrupt(0,CHECK_IR,CHANGE);
  irrecv.enableIRIn(); // Start the receiver
  inputString.reserve(200);

  //set pwm pins
  pinMode(REDPIN, OUTPUT);
  pinMode(GREENPIN, OUTPUT);
  pinMode(BLUEPIN, OUTPUT);
}

void loop() {
 serialEvent(); //call the function
  // print the string when a newline arrives:
  if (stringComplete) {
    Serial.println(inputString);
    decode_Serial(inputString);
    // clear the string:
    inputString = "";
    stringComplete = false;
  }
}

/**
 * interrupt routine for the IR receiver
 */
void CHECK_IR(){
  while(irrecv.decode(&results)){
    Serial.println(results.value, DEC);
    decode_Remote(results.value);
    irrecv.resume();
  }
}

/**
 * interprets the input from the serial interface
 */
void decode_Serial(String cmd){
  if(cmd.startsWith("color:(", 0)){
    int start = cmd.indexOf("(")+1;
    int end_r = cmd.indexOf(",");
    int end_g = cmd.indexOf(",",end_r +1);
    int end_b = cmd.indexOf(")");

    int red = cmd.substring(start,end_r).toInt();
    int green = cmd.substring(end_r+1,end_g).toInt();
    int blue = cmd.substring(end_g+1,end_b).toInt();

    setColor(red, green, blue);
  }
}

/**
 * interpretes the IR Remote codes to commands and colors
 */
void decode_Remote(unsigned long code){
    switch (code){

      //POWER ON
      case 3457774333:
        Serial.println("cmd:POWER_ON");
        setColor(255,255,255);   
        break;

      //POWER OFF
      case 3571254145:
      case 16203967:
        Serial.println("cmd:POWER_OFF");
        setColor(0,0,0);   
        break;

      //RED
      case 3898168033:
      case 16195807:
        setColor(255,0,0);  
        break;

      //GREEN
      case 2026743005:
      case 16228447:  
        setColor(0,255,0);
        break;
      
      //BLUE  
      case 2724668229:
      case 16212127:
        setColor(0,0,255);
        break;
      
      //WHITE
      case 2611188417:
      case 16244767:
        setColor(255,255,255);
        break;

      //DECREASE BRIGHTNESS
      case 16220287:
      case 3736127585:  
        decreaseBrightness();
        break;

      //INCREASE BRIGHTNESS
      case 16187647:
      case 2231595101: 
        increaseBrightness();
        break;
    }
}

/**
 * sets the color of the led stripe
 */
void setColor(unsigned char red, unsigned char green, unsigned char blue){
  //set colors
  analogWrite(REDPIN, red);
  analogWrite(GREENPIN, green);
  analogWrite(BLUEPIN, blue);

  //Save actual colors
  Red = red;
  Green = green;
  Blue = blue;

  printColor();
}

/**
 * Decrease Brightness by decreasing PWM Output by 25
 */
void decreaseBrightness(){
  if(Red >= 25){
    Red -= 25;  
  }else{
    Red = 0;
  }

  if(Green >= 25){
    Green -= 25;  
  }else{
    Green = 0;
  }

  if(Blue >= 25){
    Blue -= 25;  
  }else{
    Blue = 0;
  }
  
  setColor(Red, Green, Blue);
  printColor();
}


/**
 * Increase Brightness by increasing PWM Output by 25
 */
void increaseBrightness(){
  if(Red <= 230){
    Red += 25;  
  }else{
    Red = 255;
  }

  if(Green <= 230){
    Green += 25;  
  }else{
    Green = 255;
  }

  if(Blue <= 230){
    Blue += 25;  
  }else{
    Blue = 255;
  }
  
  setColor(Red, Green, Blue);
  printColor();
}

/*
 * Print actual color to the serial interface
 */
void printColor(){
  Serial.print("color:(");
  Serial.print(Red);
  Serial.print(",");
  Serial.print(Green);
  Serial.print(",");
  Serial.print(Blue);
  Serial.println(")");
}

/*
 * checks if the serial read is complete
 */
void serialEvent() {
  while (Serial.available()) {
    // get the new byte:
    char inChar = (char)Serial.read();
    // add it to the inputString:
    inputString += inChar;
    // if the incoming character is a newline, set a flag
    // so the main loop can do something about it:
    if (inChar == '\n') {
      stringComplete = true;
    }
  }
}
