plugins {
    id("java")
    id("com.google.cloud.tools.jib") version "3.4.4"
}

group = "com.toddljones"

repositories {
    mavenCentral()
}

jib {
    from {
        image = "eclipse-temurin:23-jre"
    }
    container {
        mainClass = "com.toddljones.wayfarer.Main"
    }
}


dependencies {
    implementation("com.google.maps:google-maps-routing:1.37.0")
    implementation("jakarta.inject:jakarta.inject-api:2.0.1")
    testImplementation("org.mockito:mockito-core:5.15.2")
    testImplementation("org.mockito:mockito-junit-jupiter:5.3.1")
    testImplementation("org.assertj:assertj-core:3.24.2")
    testImplementation(platform("org.junit:junit-bom:5.10.0"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

tasks.test {
    useJUnitPlatform()
}