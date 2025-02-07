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
        args = listOf("--config-file=/app/config.yaml")
    }
}


dependencies {
    implementation("com.google.maps:google-maps-routing:1.42.0")
    implementation("jakarta.inject:jakarta.inject-api:2.0.1")
    implementation("commons-cli:commons-cli:1.4")
    implementation("org.slf4j:slf4j-api:2.0.7")
    implementation("org.slf4j:slf4j-simple:2.0.16")
    implementation("com.fasterxml.jackson.dataformat:jackson-dataformat-yaml:2.15.0")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310:2.15.0")
    testImplementation("org.mockito:mockito-core:5.15.2")
    testImplementation("org.mockito:mockito-junit-jupiter:5.3.1")
    testImplementation("org.assertj:assertj-core:3.24.2")
    testImplementation(platform("org.junit:junit-bom:5.11.4"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    testImplementation("org.awaitility:awaitility:4.2.2")
}

tasks.test {
    useJUnitPlatform()
}