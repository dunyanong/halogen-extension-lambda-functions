const { DynamoDBClient, PutItemCommand } = require('@aws-sdk/client-dynamodb');

const dynamoDBClient = new DynamoDBClient({ region: 'ap-southeast-1' });

const DynamoDB = {
    async write(hash, fileName, tableName) {
        const timestamp = new Date().toISOString(); // Get current timestamp in ISO format
        
        const params = {
            TableName: tableName,
            Item: {
                hash: { S: hash },
                filename: { S: fileName },
                timestamp: { S: timestamp } // Add timestamp attribute
            }
        };

        try {
            const command = new PutItemCommand(params);
            const newData = await dynamoDBClient.send(command);

            if (!newData) {
                throw new Error('There was an error writing the item');
            }

            return newData;
        } catch (error) {
            console.error(error);
            throw new Error(`Error writing item: ${error.message}`);
        }
    },
};

module.exports = DynamoDB;